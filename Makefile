.PHONY: build-prod run-dev proto test

build-prod:
	export CGO_ENABLED=0
	CGO_ENABLED=0 go build -trimpath -ldflags="-s" -o build/ cmd/trader.go

run-dev:
	go run cmd/trader.go

proto:
	protoc --go_out api/bth --go-grpc_out=api/bth/ api/proto/trader.proto

test:
	go test ./...
