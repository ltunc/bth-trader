# Trader

Provides API for other internal services to execute trades.

Watches for open orders, notify other services when an order was completed or closed
Provides gRPC service for others to connect. Protobuf files are in `api/proto/`

## Env Parameters


_all parameters have prefix `BTH_`_

* `BTH_KRAKEN_API_KEY` - API key to access to Kraken API
* `BTH_KRAKEN_PRIVATE_KEY` - Private key to access to Kraken API
* `BTH_GRPC_LISTEN` - Address and port to open gRPC server on (default 0.0.0.0:5500)

## Build

    build-prod

or 

    go build cmd/trader.go