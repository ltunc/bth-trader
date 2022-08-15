package orders

import (
	"bth-trader/internal/entities"
	"reflect"
	"testing"
)

type MockObserver struct {
	Name string
}

func (m *MockObserver) Notify(_ *entities.Order) {
}

func TestStorage_Unsubscribe(t *testing.T) {
	type fields struct {
		subscribers []Observer
	}
	type args struct {
		o Observer
	}
	observers := []Observer{&MockObserver{"M1"}, &MockObserver{"M2"}, &MockObserver{"M3"}}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []Observer
	}{
		{
			name: "middle",
			fields: fields{
				subscribers: []Observer{observers[0], observers[1], observers[2]},
			},
			args: args{observers[1]},
			want: []Observer{observers[0], observers[2]},
		},
		{
			name: "last",
			fields: fields{
				subscribers: []Observer{observers[0], observers[1], observers[2]},
			},
			args: args{observers[2]},
			want: []Observer{observers[0], observers[1]},
		},
		{
			name: "first",
			fields: fields{
				subscribers: []Observer{observers[0], observers[1], observers[2]},
			},
			args: args{observers[0]},
			want: []Observer{observers[1], observers[2]},
		},
		{
			name: "not existing",
			fields: fields{
				subscribers: []Observer{observers[0], observers[1], observers[2]},
			},
			args: args{&MockObserver{"U1"}},
			want: observers,
		},
		{
			name: "empty",
			fields: fields{
				subscribers: []Observer{},
			},
			args: args{&MockObserver{"U1"}},
			want: []Observer{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Dispatcher{
				subscribers: tt.fields.subscribers,
			}
			d.Unsubscribe(tt.args.o)
			got := d.subscribers
			if !reflect.DeepEqual(tt.want, got) {
				t.Errorf("Unsubscribe() expected subscribers list %v, got %v", tt.want, got)
			}
		})
	}
}

func TestStorage_Subscribe(t *testing.T) {
	type fields struct {
		subscribers []Observer
	}
	type args struct {
		o Observer
	}
	observers := []Observer{&MockObserver{"M1"}, &MockObserver{"M2"}, &MockObserver{"M3"}}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []Observer
	}{
		{
			"empty",
			fields{subscribers: []Observer{}},
			args{observers[1]},
			[]Observer{observers[1]},
		},
		{
			"not empty",
			fields{subscribers: []Observer{observers[0], observers[1]}},
			args{observers[2]},
			observers,
		},
		{
			"existing",
			fields{subscribers: []Observer{observers[0], observers[1], observers[2]}},
			args{observers[1]},
			observers,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Dispatcher{
				subscribers: tt.fields.subscribers,
			}
			d.Subscribe(tt.args.o)
			got := d.subscribers
			if !reflect.DeepEqual(tt.want, got) {
				t.Errorf("Subscribe() expected subscribers list %#v, got %#v", tt.want, got)
			}
		})
	}
}

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
