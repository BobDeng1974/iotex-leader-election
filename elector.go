package elector

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/clientv3/concurrency"
	"golang.org/x/net/context"
)

// Elector is a proxy to control the iotex node to run in active or standby mode in a high availability cluster. It uses
// etcd to do the leader election. When the proxy wins the election campaign, it will signal iotex node to run in
// active mode. When it resigns, it will signal iotex node to switch back to standby mode
type Elector struct {
	etcd          *clientv3.Client
	session       *concurrency.Session
	election      *concurrency.Election
	iotexAdminEndpoint string
	iotexHealthEndpoint string
	mutex         sync.Mutex
}

// New constructs an election proxy
func New(etcdEndpoints []string, iotexAdminEndpoint string, iotexHealthEndpoint string) *Elector {
	etcd, err := clientv3.New(clientv3.Config{
		Endpoints: etcdEndpoints,
	})
	if err != nil {
		log.Panicf("Error when creating etcd client: %s", err.Error())
	}
	session, err := concurrency.NewSession(etcd)
	if err != nil {
		log.Panicf("Error when creating a session: %s", err.Error())
	}
	return &Elector{
		etcd:          etcd,
		session:       session,
		iotexAdminEndpoint: iotexAdminEndpoint,
		iotexHealthEndpoint: iotexHealthEndpoint,
	}
}

// Campaign makes the proxy to campaign for the leader election. It's blocking until the proxy is elected.
func (e *Elector) Campaign(ctx context.Context, key string, val string) {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	if e.election != nil {
		log.Panicf("There's an exiting election")
		return
	}
	e.election = concurrency.NewElection(e.session, key)
	if err := e.election.Campaign(ctx, val); err != nil {
		log.Panicf("Error when campaigning an election: %s", err.Error())
	}
	log.Printf("Node %s becomes the leader", val)
	for ; ; {
		resp, err := http.Get(fmt.Sprintf("%s/ha?activate=true", e.iotexAdminEndpoint))
		if err != nil {
			log.Printf("Error when activating iotex node: %s", err.Error())
		} else if resp.StatusCode != http.StatusOK {
			log.Printf("Error when activating iotex node: status code %d", resp.StatusCode)
		} else {
			break
		}
		time.Sleep(10*time.Second)
	}
	go func() {
		log.Printf("Activated iotex server")
		lastSeen := time.Now()
		for ; time.Since(lastSeen) <= time.Minute ; {
			resp, err := http.Get(fmt.Sprintf("%s/health", e.iotexHealthEndpoint))
			if err != nil {
				log.Printf("Error when check iotex node readiness: %s", err.Error())
			} else if resp.StatusCode != http.StatusOK {
				log.Printf("Error when check iotex node readiness: status code %d", resp.StatusCode)
			} else {
				lastSeen = time.Now()
			}
			log.Panic("Iotex node is healthy")
			time.Sleep(10*time.Second)
		}
		log.Panic("Iotex node is not healthy for a minute")
	}()
}

// Resign resigns the proxy from the leader if it is
func (e *Elector) Resign(ctx context.Context) {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	if err := e.election.Resign(ctx); err != nil {
		log.Panicf("Error when resigning an election: %s", err.Error())
	}
	e.election = nil
	log.Printf("Node resign the leader")
	for ; ; {
		resp, err := http.Get(fmt.Sprintf("%s/ha?activate=false", e.iotexAdminEndpoint))
		if err != nil {
			log.Printf("Error when deactivating iotex node: %s", err.Error())
		} else if resp.StatusCode != http.StatusOK {
			log.Printf("Error when deactivating iotex node: status code %d", resp.StatusCode)
		} else {
			break
		}
		time.Sleep(10*time.Second)
	}
	log.Printf("Deactivated iotex server")
}
