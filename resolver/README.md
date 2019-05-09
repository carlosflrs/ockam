# Universal Resolver Driver: did:ockam

This is a [Universal Resolver](https://github.com/decentralized-identity/universal-resolver/) driver for **did:ockam** identifiers.

## Build and Run (Docker)

```
docker build -t ockam/resolver . \
    --build-arg DISCOVERER_NAME="test.ockam.network" \
    --build-arg DISCOVERER_PORT=26657
docker run -p 8080:8080 ockam/resolver --rm
```


