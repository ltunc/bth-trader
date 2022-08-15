package kraken

import (
	"bth-trader/internal/entities"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"strconv"
	"sync"
)

const WsEndpoint = "wss://ws-auth.kraken.com"

type WsClient struct {
	token    string
	m        *sync.Mutex
	endpoint string
	conn     *websocket.Conn
	output   chan json.RawMessage
}

type SubMessage struct {
	Event        string         `json:"event,omitempty"`
	Pair         []string       `json:"pair,omitempty"`
	Subscription map[string]any `json:"subscription,omitempty"`
}

// NewWsClient creates new Websocket client
// Argument endpoint should be full address of the server to connect to
// The returned client will not be connected to the server. To do so call method *WsClient.Dial()
func NewWsClient(endpoint string) *WsClient {
	return &WsClient{
		endpoint: endpoint,
		m:        &sync.Mutex{},
	}
}

// Dial connects to remote server
// Takes address to connect from WsClient.endpoint property
// Returns original errors
func (w *WsClient) Dial() error {
	w.m.Lock()
	defer w.m.Unlock()
	if w.conn != nil {
		return nil
	}
	log.Printf("dialing kraken server: %s", w.endpoint)
	var err error
	if w.conn, _, err = websocket.DefaultDialer.Dial(w.endpoint, nil); err != nil {
		return err
	}
	return nil
}

// Subscribe sends "subscribe" event to the server
// The connection to WS server should be alive.
// Returns original error wrapped with more descriptions
func (w *WsClient) Subscribe(sub SubMessage) error {
	w.m.Lock()
	defer w.m.Unlock()
	if err := w.conn.WriteJSON(sub); err != nil {
		return fmt.Errorf("cannot send subscribe message: %w", err)
	}
	return nil
}

// Stream starts reading messages from websocket connection
// Returns a channel to which it sends all received messages.
// Can be called multiple times, but creates channel only first time,
// subsequent calls will return already opened channel.
// Closes the output channel if the websocket server returned an error.
func (w *WsClient) Stream() <-chan json.RawMessage {
	w.m.Lock()
	defer w.m.Unlock()
	if w.output != nil {
		return w.output
	}
	w.output = make(chan json.RawMessage, 100)
	go func() {
		for {
			_, msg, err := w.conn.ReadMessage()
			if err != nil {
				close(w.output)
				log.Printf("websocket read error: %v", err)
				return
			}
			w.output <- msg
		}
	}()
	return w.output
}

type AddOrderMsg struct {
	Event     string `json:"event"`
	OrderType string `json:"ordertype"`
	Pair      string `json:"pair"`
	Price     string `json:"price"`
	ReqId     int    `json:"reqid,omitempty"`
	Token     string `json:"token"`
	Type      string `json:"type"`
	UserRef   string `json:"userref"`
	Volume    string `json:"volume"`
}

func NewAddOrderMsg(refId int, pair, direction string, price, volume float64, token string) AddOrderMsg {
	msg := AddOrderMsg{
		Event:     "addOrder",
		OrderType: "limit",
		Pair:      pair,
		Price:     fmt.Sprintf("%v", price),
		ReqId:     refId,
		Token:     token,
		Type:      direction,
		UserRef:   strconv.Itoa(refId),
		Volume:    fmt.Sprintf("%v", volume),
	}
	return msg
}

func (w *WsClient) AddOrder(msg AddOrderMsg) error {
	w.m.Lock()
	defer w.m.Unlock()
	if err := w.conn.WriteJSON(msg); err != nil {
		return fmt.Errorf("cannot send addOrder message: %w", err)
	}
	return nil
}

type CancelOrderMsg struct {
	Event string   `json:"event"`
	ReqId int      `json:"reqid,omitempty"`
	TxId  []string `json:"txid"`
	Token string   `json:"token"`
}

func NewCancelOrderMsg(ord []*entities.Order, token string) CancelOrderMsg {
	var ids []string
	for _, o := range ord {
		ids = append(ids, o.OrderId)
	}
	msg := CancelOrderMsg{
		Event: "cancelOrder",
		TxId:  ids,
		Token: token,
	}
	return msg
}

func (w *WsClient) CancelOrder(msg CancelOrderMsg) error {
	w.m.Lock()
	defer w.m.Unlock()
	if err := w.conn.WriteJSON(msg); err != nil {
		return fmt.Errorf("cannot send cancelOrder message: %w", err)
	}
	return nil
}
