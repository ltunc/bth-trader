package server

import (
	"bth-trader/api/bth"
	"bth-trader/internal/entities"
	"bth-trader/internal/kraken"
	"bth-trader/internal/orders"
	"context"
	"github.com/ltunc/go-observer/observer"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"math/rand"
	"time"
)

type TraderServer struct {
	bth.UnimplementedTraderServer
	ws      *kraken.WsClient
	token   *kraken.WsAuthToken
	rnd     *rand.Rand
	od      *observer.Subject[*entities.Order]
	storage *orders.Storage
}

func NewTraderServer(ws *kraken.WsClient, token *kraken.WsAuthToken, od *observer.Subject[*entities.Order], storage *orders.Storage) *TraderServer {
	return &TraderServer{
		ws:      ws,
		token:   token,
		rnd:     rand.New(rand.NewSource(time.Now().UnixMilli())),
		od:      od,
		storage: storage,
	}
}

func (s *TraderServer) AddOrder(_ context.Context, req *bth.AddOrderRequest) (*bth.AddOrderResponse, error) {
	refId := int(s.rnd.Int31())
	orderWaiter := orders.NewWaiter(refId)
	s.od.Subscribe(orderWaiter)
	defer s.od.Unsubscribe(orderWaiter)
	msg := kraken.NewAddOrderMsg(refId, req.Pair, req.Direction, req.Price, req.Volume, s.token.Token)
	if err := s.ws.AddOrder(msg); err != nil {
		return nil, status.Errorf(codes.Internal, "cannot place an order: %v", err)
	}
	order := orderWaiter.Wait()
	if order.Status == "error" {
		return nil, status.Errorf(codes.Internal, "error when placing an order: %v", order.Error)
	}
	resp := &bth.AddOrderResponse{
		Status:  order.Status,
		RefId:   int32(refId),
		OrderId: order.OrderId,
	}
	return resp, nil
}

func (s *TraderServer) CancelOrder(_ context.Context, req *bth.CancelOrderRequest) (*bth.CancelOrderResponse, error) {
	refId := int(req.GetRefId())
	order, ok := s.storage.Find(refId)
	if !ok {
		return nil, status.Errorf(codes.NotFound, "cannot find order %v", refId)
	}
	msg := kraken.NewCancelOrderMsg([]*entities.Order{order}, s.token.Token)
	if err := s.ws.CancelOrder(msg); err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	resp := &bth.CancelOrderResponse{Status: "success"}
	return resp, nil
}

func (s *TraderServer) OrderStatus(_ context.Context, req *bth.OrderStatusRequest) (*bth.OrderStatusResponse, error) {
	refId := int(req.GetRefId())
	order, ok := s.storage.Find(refId)
	if !ok {
		return nil, status.Errorf(codes.NotFound, "cannot find order by RefId %d", refId)
	}
	resp := &bth.OrderStatusResponse{
		RefId:   int32(order.RefId),
		OrderId: order.OrderId,
		Status:  order.Status,
	}
	return resp, nil
}

// copyObs is observer that sends received order to another channel for consumption
type copyObs struct {
	ch chan *entities.Order
}

func (c *copyObs) Notify(order *entities.Order) {
	select {
	case c.ch <- order:
	default:
		log.Printf("dropped update %v, channel is full", order)
	}
}

func (s *TraderServer) StreamOrders(_ *bth.Empty, stream bth.Trader_StreamOrdersServer) error {
	inOrders := &copyObs{
		ch: make(chan *entities.Order, 100),
	}
	s.od.Subscribe(inOrders)
	defer s.od.Unsubscribe(inOrders)
	for o := range inOrders.ch {
		resp := &bth.OrderStatusResponse{
			RefId:   int32(o.RefId),
			OrderId: o.OrderId,
			Status:  o.Status,
		}
		err := stream.Send(resp)
		if err != nil {
			log.Printf("cannot send message to outgoing stream: %v", err)
			return err
		}
	}
	return nil
}
