package main

import (
	"bth-trader/api/bth"
	"bth-trader/internal/entities"
	"bth-trader/internal/kraken"
	"bth-trader/internal/kraken/decoder"
	"bth-trader/internal/orders"
	"bth-trader/internal/server"
	"bth-trader/internal/utils/env"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"os/signal"
	"time"
)

func main() {
	rest := kraken.NewRestClient(env.Get("KRAKEN_API_KEY", ""), env.Get("KRAKEN_PRIVATE_KEY", ""))
	token, err := rest.WsToken()
	if err != nil {
		log.Printf("cannot receive auth token for Websocket requests: %v", err)
		return
	}
	ws := kraken.NewWsClient(kraken.WsEndpoint)
	if err := ws.Dial(); err != nil {
		log.Fatalf("cannot dial kraken: %v", err)
	}
	if err := subKraken(ws, token); err != nil {
		log.Fatalf("cannot configure kraken WS: %v", err)
	}
	// read stream of messages from the server
	out := &decoder.Outputs{
		Orders: make(chan *entities.Order, 50),
		Trades: make(chan *entities.Trade, 50),
	}
	go decoder.DecodeStream(ws.Stream(), out)
	go consume(out)
	od := &orders.Dispatcher{}
	storage := orders.NewStorage()
	go runStorageGc(storage)
	od.Subscribe(storage)
	go od.ReadFrom(out.Orders)
	lis, err := net.Listen("tcp", env.Get("GRPC_LISTEN", "127.0.0.1:5500"))
	if err != nil {
		log.Fatalf("cannot open port: %v", err)
	}
	go runGrpc(lis, ws, token, od, storage)
	wait()
}

// subKraken subscribes kraken WS client for all necessary channels
func subKraken(ws *kraken.WsClient, token *kraken.WsAuthToken) error {
	openOrders := kraken.SubMessage{
		Event:        "subscribe",
		Subscription: map[string]any{"name": "openOrders", "token": token.Token},
	}
	if err := ws.Subscribe(openOrders); err != nil {
		return fmt.Errorf("cannot subscribe to open orders: %w", err)
	}
	ownTrades := kraken.SubMessage{
		Event:        "subscribe",
		Subscription: map[string]any{"name": "ownTrades", "token": token.Token},
	}
	if err := ws.Subscribe(ownTrades); err != nil {
		return fmt.Errorf("cannot subscribe to own trades: %w", err)
	}
	return nil
}

// runStorageGc executes cleaning process of the storage, removes old closed/finished orders
func runStorageGc(s *orders.Storage) {
	ticker := time.NewTicker(time.Second * 2)
	for range ticker.C {
		orders.Cleanup(s)
	}
}

// runGrpc prepares and starts gRPC server
func runGrpc(lis net.Listener, ws *kraken.WsClient, t *kraken.WsAuthToken, dispatcher *orders.Dispatcher, storage *orders.Storage) {
	var opts []grpc.ServerOption
	srv := grpc.NewServer(opts...)
	bth.RegisterTraderServer(srv, server.NewTraderServer(ws, t, dispatcher, storage))
	log.Fatal(srv.Serve(lis))
}

// consume reads messages from output and prints it to logs
// it is a temporary/debug reader
func consume(out *decoder.Outputs) {
	for {
		select {
		//case order := <-out.Orders:
		//	log.Printf("order: %v", order)
		case trade := <-out.Trades:
			log.Printf("trade: %v", trade)
		}
	}
}

// wait blocks goroutine until SIGINT received
func wait() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	log.Printf("interrupted")
}
