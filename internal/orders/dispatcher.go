package orders

import (
	"bth-trader/internal/entities"
	"github.com/ltunc/go-observer/observer"
)

func NewDispatcher() *observer.Subject[*entities.Order] {
	return &observer.Subject[*entities.Order]{}
}

// ReadFrom reads orders from the channel and fires events in the dispatcher
// notifies all observers about new update
func ReadFrom(dispatcher *observer.Subject[*entities.Order], input <-chan *entities.Order) {
	for order := range input {
		dispatcher.Fire(order)
	}
}

// Waiter is an observer that waits for concrete order
type Waiter struct {
	expectRefId int
	results     chan *entities.Order
}

// NewWaiter creates an observer that blocks a goroutine until update for an order with refId received
func NewWaiter(refId int) *Waiter {
	return &Waiter{
		expectRefId: refId,
		results:     make(chan *entities.Order, 1),
	}
}

func (w *Waiter) Notify(o *entities.Order) {
	if o.RefId == w.expectRefId {
		select {
		case w.results <- o:
		default:
			//log.Printf("cannot send the order to results of a waiter, channel is full")
		}
	}
}

// Wait blocks execution of goroutine until an order appears
func (w *Waiter) Wait() *entities.Order {
	order := <-w.results
	return order
}
