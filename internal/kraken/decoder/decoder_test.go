package decoder

import (
	"bth-trader/internal/entities"
	"encoding/json"
	"io"
	"log"
	"os"
	"reflect"
	"testing"
)

//func TestMain(m *testing.M) {
//	// suppress logging output during tests
//	log.SetOutput(io.Discard)
//	os.Exit(m.Run())
//}

func TestDecodeStream(t *testing.T) {
	type args struct {
		in  chan json.RawMessage
		out *Outputs
	}
	type testOutput struct {
		orders []*entities.Order
		trades []*entities.Trade
	}
	tests := []struct {
		name       string
		args       args
		inMessages []json.RawMessage
		wantOut    testOutput
	}{
		{
			name:       "system status",
			inMessages: []json.RawMessage{json.RawMessage(`{"connectionID":16569497294059334297,"event":"systemStatus","status":"online","version":"1.9.0"}`)},
			args: args{
				make(chan json.RawMessage, 6),
				&Outputs{Orders: make(chan *entities.Order, 100), Trades: make(chan *entities.Trade, 100)},
			},
			wantOut: testOutput{},
		},
		{
			name:       "heartbeat",
			inMessages: []json.RawMessage{json.RawMessage(`{"event":"heartbeat"}`)},
			args: args{
				make(chan json.RawMessage, 6),
				&Outputs{Orders: make(chan *entities.Order, 100), Trades: make(chan *entities.Trade, 100)},
			},
			wantOut: testOutput{},
		},
		{
			name:       "addOrder success",
			inMessages: []json.RawMessage{json.RawMessage(`{"event":"addOrderStatus", "reqid": 112233, "status": "ok", "txid": "ABCDEF-ABCD2-ABCDE3", "descr": "test order"}`)},
			args: args{
				make(chan json.RawMessage, 6),
				&Outputs{Orders: make(chan *entities.Order, 100), Trades: make(chan *entities.Trade, 100)},
			},
			wantOut: testOutput{orders: []*entities.Order{
				{OrderId: "ABCDEF-ABCD2-ABCDE3", RefId: 112233, Status: "open"},
			}},
		},
		{
			name:       "addOrder error",
			inMessages: []json.RawMessage{json.RawMessage(`{"event":"addOrderStatusStatus", "reqid": 223344, "status": "error", "errorMessage": "TestErrorMsg"}`)},
			args: args{
				make(chan json.RawMessage, 6),
				&Outputs{Orders: make(chan *entities.Order, 100), Trades: make(chan *entities.Trade, 100)},
			},
			wantOut: testOutput{orders: []*entities.Order{
				{RefId: 223344, Status: "error", Error: "TestErrorMsg"},
			}},
		},
		{
			name:       "order open",
			inMessages: []json.RawMessage{json.RawMessage(`[[{"ABCDEF-ABCD2-ABCDE3":{"avg_price":"0.00000","cost":"0.00000","descr":{"close":null,"leverage":null,"order":"buy 0.90101951 XBT/EUR @ limit 23302.00000","ordertype":"limit","pair":"XBT/EUR","price":"23302.00000","price2":"0.00000","type":"buy"},"expiretm":null,"fee":"0.00000","limitprice":"0.00000","misc":"","oflags":"fciq","opentm":"1660000011.012345","refid":123456,"starttm":null,"status":"open","stopprice":"0.00000","timeinforce":"GTC","userref":123456,"vol":"0.90101951","vol_exec":"0.00000000"}}],"openOrders",{"sequence":1}]`)},
			args: args{
				make(chan json.RawMessage, 6),
				&Outputs{Orders: make(chan *entities.Order, 100), Trades: make(chan *entities.Trade, 100)},
			},
			wantOut: testOutput{orders: []*entities.Order{
				{OrderId: "ABCDEF-ABCD2-ABCDE3", RefId: 123456, Status: "open"},
			}},
		},
		{
			name:       "order update",
			inMessages: []json.RawMessage{json.RawMessage(`[[{"ABCDEF-ABCD2-ABCDE4":{"lastupdated":"1650000019.012345","status":"canceled","vol_exec":"0.00000000","cost":"0.00000","fee":"0.00000","avg_price":"0.00000","userref":778899,"cancel_reason":"User requested"}}],"openOrders",{"sequence":2}]`)},
			args: args{
				make(chan json.RawMessage, 6),
				&Outputs{Orders: make(chan *entities.Order, 100), Trades: make(chan *entities.Trade, 100)},
			},
			wantOut: testOutput{orders: []*entities.Order{
				{OrderId: "ABCDEF-ABCD2-ABCDE4", RefId: 778899, Status: "canceled"},
			}},
		},
		{
			name:       "trade",
			inMessages: []json.RawMessage{json.RawMessage(`[[{"TTTTTT-AAAA1-EEEEE1":{"cost":"100.14230","fee":"0.16023","margin":"0.00000","ordertxid":"OZXDAA-A10A1-0ABCDE","ordertype":"limit","pair":"ETH/EUR","postxid":"TABCDE-ABCD1-ABCDE2","price":"1728.40000","time":"1650000011.061588","type":"sell","vol":"0.05793931"}}],"ownTrades",{"sequence":1}]`)},
			args: args{
				make(chan json.RawMessage, 6),
				&Outputs{Orders: make(chan *entities.Order, 100), Trades: make(chan *entities.Trade, 100)},
			},
			wantOut: testOutput{},
		},
	}
	log.SetOutput(io.Discard)
	defer log.SetOutput(os.Stderr)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, m := range tt.inMessages {
				tt.args.in <- m
			}
			close(tt.args.in)
			DecodeStream(tt.args.in, tt.args.out)
			close(tt.args.out.Orders)
			close(tt.args.out.Trades)
			var gotOrders []*entities.Order
			for o := range tt.args.out.Orders {
				gotOrders = append(gotOrders, o)
			}
			var gotTrades []*entities.Trade
			for t := range tt.args.out.Trades {
				gotTrades = append(gotTrades, t)
			}
			if !reflect.DeepEqual(gotTrades, tt.wantOut.trades) {
				t.Errorf("DecodeStream() expected output.Trades = %v, got %v", tt.wantOut.trades, gotTrades)
				for k, o := range tt.wantOut.trades {
					t.Logf("#%d: want: %v", k, o)
					t.Logf("#%d: got: %v", k, gotTrades[k])
				}
			}
			if !reflect.DeepEqual(gotOrders, tt.wantOut.orders) {
				t.Errorf("DecodeStream() expected output.Orders = %v, got %v", tt.wantOut.orders, gotOrders)
				for k, o := range tt.wantOut.orders {
					t.Logf("#%d: want: %v", k, o)
					t.Logf("#%d: got: %v", k, gotOrders[k])
				}
			}
		})
	}
}

func Test_detectType(t *testing.T) {
	type args struct {
		rawData any
	}
	tests := []struct {
		name string
		args args
		want msgType
	}{
		{
			name: "heartbeat",
			args: args{map[string]any{"event": "heartbeat"}},
			want: msgHeartbeat,
		},
		{
			name: "system-status",
			args: args{map[string]any{"connectionID": "101", "event": "systemStatus", "status": "online", "version": "1.9.0"}},
			want: msgSysStatus,
		},
		{
			name: "sub-status",
			args: args{map[string]any{"channelName": "openOrders", "event": "subscriptionStatus", "status": "subscribed", "subscription": map[string]any{"maxratecount": 125, "name": "openOrders"}}},
			want: msgSubStatus,
		},
		{
			name: "add-order",
			args: args{map[string]any{"channelName": "addOrder", "event": "addOrderStatus", "status": "success"}},
			want: msgAddOrderStatus,
		},
		{
			name: "add-order",
			args: args{map[string]any{"channelName": "addOrder", "event": "addOrderStatusStatus", "status": "error"}},
			want: msgAddOrderStatus,
		},
		{
			name: "cancel-order",
			args: args{map[string]any{"channelName": "cancelOrder", "event": "cancelOrderStatus", "status": "ok"}},
			want: msgCancelOrderStatus,
		},
		{
			name: "cancel-order error",
			args: args{map[string]any{"channelName": "cancelOrder", "event": "cancelOrderStatus", "status": "error", "errorMessage": "test error"}},
			want: msgCancelOrderStatus,
		},
		{
			name: "order",
			args: args{[]any{[]any{map[string]any{"1ABCDE-FGHIJ-12345A": map[string]any{"avg_price": "0.00000", "cost": "0.00000", "descr": map[string]interface{}{"close": interface{}(nil), "leverage": interface{}(nil), "order": "buy 0.90011223 XBT/EUR @ limit 23302.00000", "ordertype": "limit", "pair": "XBT/EUR", "price": "23302.00000", "price2": "0.00000", "type": "buy"}, "expiretm": interface{}(nil), "fee": "0.00000", "limitprice": "0.00000", "misc": "", "oflags": "fciq", "opentm": "1650000011.012345", "refid": 123456, "starttm": interface{}(nil), "status": "open", "stopprice": "0.00000", "timeinforce": "GTC", "userref": "123456", "vol": "0.90011223", "vol_exec": "0.00000000"}}}, "openOrders", map[string]any{"sequence": "1"}}},
			want: msgOrder,
		},
		{
			name: "order-status",
			args: args{[]any{[]any{map[string]any{"1ABCDE-FGHIJ-12345A": map[string]any{"lastupdated": "1650000019.012345", "status": "canceled", "vol_exec": "0.00000000", "cost": "0.00000", "fee": "0.00000", "avg_price": "0.00000", "userref": 123456, "cancel_reason": "User requested"}}}, "openOrders", map[string]any{"sequence": "2"}}},
			want: msgOrder,
		},
		{
			name: "trade",
			args: args{[]any{[]any{map[string]any{"TTTTTT-AAAA1-EEEEE1": map[string]any{"cost": "100.14230", "fee": "0.16023", "margin": "0.00000", "ordertxid": "OZXDAA-A10A1-0ABCDE", "ordertype": "limit", "pair": "ETH/EUR", "postxid": "TABCDE-ABCD1-ABCDE2", "price": "1728.40000", "time": "1650000011.061588", "type": "sell", "vol": "0.05793931"}}}, "ownTrades", map[string]any{"sequence": "1"}}},
			want: msgTrade,
		},
		{
			name: "unknown",
			args: args{map[string]any{"channelName": "something", "event": "unexpected", "key": "value"}},
			want: msgUnknown,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := detectType(tt.args.rawData); got != tt.want {
				t.Errorf("detectType() = %v, want %v", got, tt.want)
			}
		})
	}
}
