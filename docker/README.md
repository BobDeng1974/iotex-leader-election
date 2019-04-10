# Docker Compose Example

- Update `config.yaml` for iotex server by adding the following:

```
network:
  ...
  masterKey: {producerPrivKey-haId}

system:
  active: false
```