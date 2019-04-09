package elector

import (
	"fmt"
	"log"
	"net/http"
	"sync"

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
	iotexEndpoint string
	mutex         sync.Mutex
}

// New constructs an election proxy
func New(etcdEndpoints []string, iotexEndpoint string) *Elector {
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
		iotexEndpoint: iotexEndpoint,
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
	resp, err := http.Get(fmt.Sprintf("%s/ha?activate=true", e.iotexEndpoint))
	if err != nil {
		log.Panicf("Error when activating iotex node: %s", err.Error())
	}
	if resp.StatusCode != http.StatusOK {
		log.Panicf("Error when activating iotex node: status code %d", resp.StatusCode)
	}
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
	resp, err := http.Get(fmt.Sprintf("%s/ha?activate=false", e.iotexEndpoint))
	if err != nil {
		log.Panicf("Error when deactivating iotex node: %s", err.Error())
	}
	if resp.StatusCode != http.StatusOK {
		log.Panicf("Error when deactivating iotex node: status code %d", resp.StatusCode)
	}
}
