# IoTeX Leader Election

This is a proxy to use [etcd](https://github.com/etcd-io/etcd) to do leader election, and control
[iotex node](https://github.com/iotexproject/iotex-core/) to run in either active or standby mode in a high availability
(HA) cluster.

In the HA cluster, all nodes are using the same key, but only one node will actively participant into consensus, while
the others will only sync blocks. 

## Usage

- Local:

```
make build
bin/elector -etcd=[etcd endpoint] -iotex=[iotex endpoint] -key=[e.g., /iotex-server] -value=[identifier in the cluster]
```

- Docker:

```
docker run -d zjshen/iotex-leader-election:latest elector -etcd=[etcd endpoint] -iotex=[iotex endpoint] -key=[e.g., /iotex-server] -value=[identifier in the cluster]
```

## Examples

- [docker-compose](docker/README.md)