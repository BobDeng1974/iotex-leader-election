package main

import (
	"context"
	"flag"

	"github.com/zjshen14/iotex-leader-election"
)

var (
	etcd  string
	iotexAdmin string
	iotexProbe string
	key   string
	val   string
)

func init() {
	flag.StringVar(&etcd, "etcd", "http://127.0.0.1:2379", "Endpoint of etcd")
	flag.StringVar(&iotexAdmin, "iotexAdmin", "http://127.0.0.1:9009", "Endpoint of iotex admin")
	flag.StringVar(&iotexProbe, "iotexProbe", "http://127.0.0.1:8080", "Endpoint of iotex probe")
	flag.StringVar(&key, "key", "/iotex-server", "Key to store in etcd")
	flag.StringVar(&val, "value", "default", "Value to store in etcd")
	flag.Parse()
}

func main() {
	e := elector.New([]string{etcd}, iotexAdmin, iotexProbe)
	ctx, _ := context.WithCancel(context.TODO())
	e.Campaign(ctx, key, val)
	defer e.Resign(context.Background())

	select {}
}
