package orders

import (
	"bth-trader/internal/entities"
	"sync"
)

// Dispatcher reads from order channel and notifies all subscribed observers about new update
type Dispatcher struct {
	// XXX: sync.Map may be more useful
	subscribers []Observer
	sync.Mutex
}

// ReadFrom reads orders from the channel into the storage
// notifies all observers about new update
func (d *Dispatcher) ReadFrom(inp <-chan *entities.Order) {
	for order := range inp {
		d.notify(order)
	}
}

func (d *Dispatcher) notify(order *entities.Order) {
	d.Lock()
	defer d.Unlock()
	for _, o := range d.subscribers {
		o.Notify(order)
	}
}

// Subscribe add an object to the list of notifications receivers
// who will be notified when new order or update appears in the database
func (d *Dispatcher) Subscribe(o Observer) {
	d.Lock()
	defer d.Unlock()
	for _, v := range d.subscribers {
		if v == o {
			return
		}
	}
	d.subscribers = append(d.subscribers, o)
}

// Unsubscribe removes an observer from the list of objects subscribed to orders
func (d *Dispatcher) Unsubscribe(o Observer) {
	d.Lock()
	defer d.Unlock()
	for k, v := range d.subscribers {
		if v == o {
			d.subscribers = append(d.subscribers[:k], d.subscribers[k+1:]...)
			return
		}
	}
}

type Observer interface {
	Notify(*entities.Order)
}

// Waiter is an observer that waits for concrete order
type Waiter struct {
	expectRefId int
	results     chan *entities.Order
}

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

//// WaitFor blocks execution of goroutine until order with refId appears in the storage
//// if the order already in the storage returns it and wake up waiting goroutine
//func WaitFor(refId int, storage *Storage) *entities.Order {
//	wtr := NewWaiter(refId)
//	storage.Subscribe(wtr)
//	return wtr.Wait()
//}
