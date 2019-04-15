# IoTeX Leader Election

[IoTeX server](https://github.com/iotexproject/iotex-core/) now supports running in active or standby mode. Only when running in active mode, the node will participate into consensus and produce blocks, otherwise, it will only passively listen to the blocks. This could be controlled by using the adming endpoint (e.g. `http://localhost:9009/ha`). Given this, it's easier to setup a backup node, which could use the same operator key. You only need to manually make one node run in active mode and the others run in standby mode.

A more advanced way is use leader election service to do the aforementioned orchestration automatically. This repo provides a proxy to use [etcd](https://github.com/etcd-io/etcd) to do leader election, and control iotex server to run in either active or standby mode in a high availability (HA) cluster.

In the HA cluster, all nodes are using the same key, but only one node will actively participant into consensus, while
the others will only sync blocks. 

## Usage

- Local:

```
make build
bin/elector -etcd=[etcd endpoint] -iotexAdmin=[iotex admin endpoint] iotexProbe=[iotex health endpoint] -key=[e.g., /iotex-server] -value=[identifier in the cluster]
```

- Docker:

```
docker run -d zjshen/iotex-leader-election:latest elector -etcd=[etcd endpoint] -iotexAdmin=[iotex admin endpoint] iotexProbe=[iotex health endpoint] -key=[e.g., /iotex-server] -value=[identifier in the cluster]
```

## Examples

Here are the examples about how to use elector together with iotex server to deploy the HA delegate cluster:

- [docker-compose](docker/README.md)

## Dependency

The leader election solution here depends on [etcd](https://github.com/etcd-io/etcd), you need to deploy it to
orchestrate the HA delegate cluster.

I deployed etcd via [helm chart](https://github.com/bitnami/charts/tree/master/bitnami/etcd), but there are many other
way to deploy it.

## Misc

IoTeX's HA feature is orthogonal to the leader election service. You can build similar solution based on
[Kubernetes election](https://github.com/kubernetes/contrib/tree/master/election),
[Zookeeper](https://zookeeper.apache.org/doc/current/recipes.html#sc_leaderElection) and etc.
