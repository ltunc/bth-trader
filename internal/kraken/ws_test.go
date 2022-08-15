package kraken

import (
	"reflect"
	"testing"
)

func TestNewAddOrderMsg(t *testing.T) {
	type args struct {
		refId     int
		pair      string
		direction string
		price     float64
		volume    float64
		token     string
	}
	tests := []struct {
		name string
		args args
		want AddOrderMsg
	}{
		{
			name: "basic",
			args: args{
				refId:     1234567,
				pair:      "XBT/USD",
				direction: "sell",
				price:     22110.19,
				volume:    0.00235101,
				token:     "some-token",
			},
			want: AddOrderMsg{
				Event:     "addOrder",
				OrderType: "limit",
				Pair:      "XBT/USD",
				Price:     "22110.19",
				ReqId:     1234567,
				Token:     "some-token",
				Type:      "sell",
				UserRef:   "1234567",
				Volume:    "0.00235101",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewAddOrderMsg(tt.args.refId, tt.args.pair, tt.args.direction, tt.args.price, tt.args.volume, tt.args.token); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewAddOrderMsg() = %v, want %v", got, tt.want)
			}
		})
	}
}
