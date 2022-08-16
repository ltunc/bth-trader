package orders

import (
	"bth-trader/internal/entities"
	"github.com/ltunc/go-observer/observer"
	"reflect"
	"testing"
)

func TestWaiter_Wait(t *testing.T) {
	type fields struct {
		expectRefId int
		results     chan *entities.Order
	}
	tests := []struct {
		name          string
		fields        fields
		notifications []*entities.Order
		want          *entities.Order
	}{
		{
			name: "basic",
			fields: fields{
				expectRefId: 11,
			},
			notifications: []*entities.Order{{RefId: 11, OrderId: "test 1"}},
			want:          &entities.Order{RefId: 11, OrderId: "test 1"},
		},
		{
			name: "repeated",
			fields: fields{
				expectRefId: 11,
			},
			notifications: []*entities.Order{{RefId: 11, OrderId: "test 1"}, {RefId: 11, OrderId: "test 1"}},
			want:          &entities.Order{RefId: 11, OrderId: "test 1"},
		},
		{
			name: "many messages",
			fields: fields{
				expectRefId: 6,
			},
			notifications: []*entities.Order{{RefId: 5, OrderId: "test 5"}, {RefId: 6, OrderId: "test 6"}, {RefId: 7, OrderId: "test 7"}},
			want:          &entities.Order{RefId: 6, OrderId: "test 6"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := NewWaiter(tt.fields.expectRefId)
			for _, n := range tt.notifications {
				w.Notify(n)
			}
			if got := w.Wait(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Wait() = %v, want %v", got, tt.want)
			}
		})
	}
}

type mockObserver struct {
	calls []*entities.Order
}

func (m *mockObserver) Notify(o *entities.Order) {
	m.calls = append(m.calls, o)
}

func TestReadFrom(t *testing.T) {
	type args struct {
		dispatcher *observer.Subject[*entities.Order]
		input      chan *entities.Order
	}
	orders := []*entities.Order{{OrderId: "AA10", RefId: 10}, {OrderId: "AA11", RefId: 11}}
	tests := []struct {
		name       string
		args       args
		sendOrders []*entities.Order
	}{
		{
			"basic",
			args{
				dispatcher: &observer.Subject[*entities.Order]{},
				input:      make(chan *entities.Order, 10),
			},
			[]*entities.Order{orders[0], orders[1]},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, o := range tt.sendOrders {
				tt.args.input <- o
			}
			close(tt.args.input)
			obs := &mockObserver{}
			tt.args.dispatcher.Subscribe(obs)
			ReadFrom(tt.args.dispatcher, tt.args.input)
			if got := obs.calls; !reflect.DeepEqual(got, tt.sendOrders) {
				t.Errorf("ReadFrom() got calls %v, want %v", got, tt.sendOrders)
			}
		})
	}
}
