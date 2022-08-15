FROM golang:alpine as builder
WORKDIR /code
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o build/trader cmd/trader.go

FROM scratch
WORKDIR /app
COPY --from=builder /code/build/* /usr/bin/
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
ENV BTH_GRPC_LISTEN="0.0.0.0:5500" \
    BTH_KRAKEN_API_KEY="" \
    BTH_KRAKEN_PRIVATE_KEY=""
EXPOSE 5500
ENTRYPOINT ["trader"]
