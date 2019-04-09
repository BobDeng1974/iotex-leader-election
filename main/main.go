package main

import (
	"context"
	"flag"
	"log"

	elector "github.com/zjshen14/iotex-leader-election"
)

var (
	etcd  string
	iotex string
	key   string
	val   string
)

func init() {
	flag.StringVar(&etcd, "etcd", "http://127.0.0.1:2379", "Endpoint of etcd")
	flag.StringVar(&iotex, "iotex", "http://127.0.0.1:9009", "Endpoint of iotex")
	flag.StringVar(&key, "key", "/iotex-server", "Key to store in etcd")
	flag.StringVar(&val, "value", "default", "Value to store in etcd")
	flag.Parse()
}

func main() {
	e := elector.New([]string{etcd}, iotex)
	e.Campaign(context.Background(), key, val)
	log.Print("ok")
	defer e.Resign(context.Background())

	select {}
}
