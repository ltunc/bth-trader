package entities

import "time"

type Order struct {
	OrderId string
	RefId   int
	// Currently we care only about order status
	Status string
	Error  string
}

type Balances map[string]float64

type Trade struct {
	Cost       float64
	Fee        float64
	Margin     float64
	OrderId    string
	OrderType  string
	Pair       string
	PositionId string
	Price      float64
	RefId      int
	Time       time.Time
	Type       string
	Volume     float64
}
