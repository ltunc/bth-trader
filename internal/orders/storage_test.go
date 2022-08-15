package orders

import (
	"bth-trader/internal/entities"
	"reflect"
	"sync"
	"testing"
	"time"
)

func TestStorage_Add(t *testing.T) {
	type fields struct {
		buffer map[int]*entities.Order
	}
	type args struct {
		order *entities.Order
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   map[int]*entities.Order
	}{
		{
			"to empty",
			fields{buffer: map[int]*entities.Order{}},
			args{&entities.Order{OrderId: "ABC", RefId: 100}},
			map[int]*entities.Order{100: {OrderId: "ABC", RefId: 100}},
		},
		{
			"not empty",
			fields{buffer: map[int]*entities.Order{
				103: {OrderId: "ABC103", RefId: 103},
				90:  {OrderId: "ABC090", RefId: 90},
			}},
			args{&entities.Order{OrderId: "ABC101", RefId: 101}},
			map[int]*entities.Order{
				103: {OrderId: "ABC103", RefId: 103},
				90:  {OrderId: "ABC090", RefId: 90},
				101: {OrderId: "ABC101", RefId: 101},
			},
		},
		{
			"no refId",
			fields{buffer: map[int]*entities.Order{
				103: {OrderId: "ABC103", RefId: 103},
				90:  {OrderId: "ABC090", RefId: 90},
			}},
			args{&entities.Order{OrderId: "ABC101"}},
			map[int]*entities.Order{
				103: {OrderId: "ABC103", RefId: 103},
				90:  {OrderId: "ABC090", RefId: 90},
			},
		},
		{
			"existing",
			fields{buffer: map[int]*entities.Order{
				91:  {OrderId: "ABC091", RefId: 91},
				102: {OrderId: "ABC102", RefId: 102, Status: "pending"},
			}},
			args{&entities.Order{OrderId: "ABC102", RefId: 102, Status: "open"}},
			map[int]*entities.Order{
				91:  {OrderId: "ABC091", RefId: 91},
				102: {OrderId: "ABC102", RefId: 102, Status: "open"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Storage{buffer: tt.fields.buffer, mu: &sync.Mutex{}}
			s.Add(tt.args.order)
			if got := s.buffer; !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Add() got buffer %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStorage_Find(t *testing.T) {
	type fields struct {
		buffer map[int]*entities.Order
	}
	type args struct {
		refId int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *entities.Order
		want1  bool
	}{
		{
			"empty",
			fields{map[int]*entities.Order{}},
			args{refId: 777},
			nil,
			false,
		},
		{
			"found",
			fields{map[int]*entities.Order{
				201: {RefId: 201, OrderId: "ABC201"},
				199: {RefId: 199, OrderId: "ABC199"},
				300: {RefId: 300, OrderId: "ABC300"},
			}},
			args{refId: 199},
			&entities.Order{RefId: 199, OrderId: "ABC199"},
			true,
		},
		{
			"notfound",
			fields{map[int]*entities.Order{
				201: {RefId: 201, OrderId: "ABC201"},
				300: {RefId: 300, OrderId: "ABC300"},
			}},
			args{refId: 199},
			nil,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Storage{buffer: tt.fields.buffer, mu: &sync.Mutex{}}
			got, got1 := s.Find(tt.args.refId)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Find() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Find() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestStorage_Notify(t *testing.T) {
	type fields struct {
		buffer map[int]*entities.Order
	}
	type args struct {
		order *entities.Order
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   map[int]*entities.Order
	}{
		{
			"basic",
			fields{map[int]*entities.Order{}},
			args{&entities.Order{OrderId: "AABB123", RefId: 123}},
			map[int]*entities.Order{123: {OrderId: "AABB123", RefId: 123}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Storage{buffer: tt.fields.buffer, mu: &sync.Mutex{}}
			s.Notify(tt.args.order)
			if got := s.buffer; !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Find() got buffer %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStorage_Remove(t *testing.T) {
	type fields struct {
		buffer map[int]*entities.Order
	}
	type args struct {
		refId int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   map[int]*entities.Order
	}{
		{
			"empty",
			fields{map[int]*entities.Order{}},
			args{refId: 555},
			map[int]*entities.Order{},
		},
		{
			"existing",
			fields{map[int]*entities.Order{
				152: {RefId: 152, OrderId: "ABC152"},
				321: {RefId: 321, OrderId: "ABC321"},
				400: {RefId: 400, OrderId: "ABC400"},
			}},
			args{refId: 321},
			map[int]*entities.Order{
				152: {RefId: 152, OrderId: "ABC152"},
				400: {RefId: 400, OrderId: "ABC400"},
			},
		},
		{
			"not existing",
			fields{map[int]*entities.Order{
				152: {RefId: 152, OrderId: "ABC152"},
				400: {RefId: 400, OrderId: "ABC400"},
			}},
			args{refId: 200},
			map[int]*entities.Order{
				152: {RefId: 152, OrderId: "ABC152"},
				400: {RefId: 400, OrderId: "ABC400"},
			},
		},
		{
			"last",
			fields{map[int]*entities.Order{
				152: {RefId: 152, OrderId: "ABC152"},
			}},
			args{refId: 152},
			map[int]*entities.Order{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Storage{buffer: tt.fields.buffer, mu: &sync.Mutex{}}
			s.Remove(tt.args.refId)
			if got := s.buffer; !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Remove() got buffer %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStorage_ByOrderId(t *testing.T) {
	type fields struct {
		buffer map[int]*entities.Order
	}
	type args struct {
		orderId string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *entities.Order
		want1  bool
	}{
		{
			"found",
			fields{map[int]*entities.Order{
				201: {RefId: 201, OrderId: "ABC201"},
				199: {RefId: 199, OrderId: "ABC199"},
				300: {RefId: 300, OrderId: "ABC300"},
			}},
			args{"ABC199"},
			&entities.Order{RefId: 199, OrderId: "ABC199"},
			true,
		},
		{
			"not found",
			fields{map[int]*entities.Order{
				201: {RefId: 201, OrderId: "ABC201"},
				199: {RefId: 199, OrderId: "ABC199"},
			}},
			args{"ABC100"},
			nil,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Storage{buffer: tt.fields.buffer, mu: &sync.Mutex{}}
			got, got1 := s.ByOrderId(tt.args.orderId)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ByOrderId() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("ByOrderId() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestCleanup(t *testing.T) {
	type args struct {
		s *Storage
	}
	now := time.Now()
	tests := []struct {
		name    string
		args    args
		want    map[int]*entities.Order
		wantDel map[int]time.Time
	}{
		{
			"full",
			args{&Storage{
				buffer: map[int]*entities.Order{
					9:  {OrderId: "ABC010", RefId: 9, Status: "closed"},
					10: {OrderId: "ABC010", RefId: 10, Status: "pending"},
					5:  {OrderId: "ABC010", RefId: 5, Status: "canceled"},
					7:  {OrderId: "ABC010", RefId: 7, Status: "open"},
				},
				deleteAt: map[int]time.Time{
					9: now,
				},
				mu: &sync.Mutex{},
			}},
			map[int]*entities.Order{
				10: {OrderId: "ABC010", RefId: 10, Status: "pending"},
				5:  {OrderId: "ABC010", RefId: 5, Status: "canceled"},
				7:  {OrderId: "ABC010", RefId: 7, Status: "open"},
			},
			map[int]time.Time{
				5: now.Add(cancelTtl),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Cleanup(tt.args.s)
			got := tt.args.s.buffer
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Cleanup() got buffer %v, want %v", got, tt.want)
			}
			//gotDel := tt.args.s.deleteAt
			//if !reflect.DeepEqual(gotDel, tt.wantDel) {
			//	t.Errorf("Cleanup() got deleteAt %v, want %v", gotDel, tt.wantDel)
			//}
		})
	}
}
