package orders

import (
	"bth-trader/internal/entities"
	"log"
	"sync"
	"time"
)

// cancelTtl is time that canceled orders should live in the storage
// after that time canceled orders removed from the storage
const cancelTtl time.Duration = time.Second * 60

// Storage stores orders in the memory and provides access to them
// implements Observer interface, so it can be subscribed to new orders from the Dispatcher
type Storage struct {
	buffer   map[int]*entities.Order
	deleteAt map[int]time.Time
	mu       *sync.Mutex
}

func Cleanup(s *Storage) {
	s.mu.Lock()
	defer s.mu.Unlock()
	now := time.Now()
	for k, o := range s.buffer {
		// ignore orders in progress
		if o.Status == "pending" || o.Status == "open" || o.Status == "opened" {
			continue
		}
		if dt, ok := s.deleteAt[k]; ok {
			if dt.Sub(now) < 0 {
				delete(s.buffer, k)
				delete(s.deleteAt, k)
			}
		} else {
			s.deleteAt[k] = now.Add(cancelTtl)
		}
	}
}

// NewStorage creates new Storage object ready to store orders
func NewStorage() *Storage {
	return &Storage{
		buffer:   make(map[int]*entities.Order),
		mu:       &sync.Mutex{},
		deleteAt: make(map[int]time.Time),
	}
}

// Notify notifies the storage about new order or an update for it
func (s *Storage) Notify(order *entities.Order) {
	s.Add(order)
}

// Add adds an order to the storage
func (s *Storage) Add(order *entities.Order) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if order.RefId == 0 {
		log.Printf("an order without RefId, ignore: %v", order)
		// store only orders with refId
		return
	}
	s.buffer[order.RefId] = order
}

// Remove removes an order from the storage
func (s *Storage) Remove(refId int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.buffer[refId]; ok {
		delete(s.buffer, refId)
	}
}

// Find search an order by its refId in the storage and returns it if found
// returns false as second argument if the order was not found
func (s *Storage) Find(refId int) (*entities.Order, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if order, ok := s.buffer[refId]; ok {
		return order, ok
	}
	return nil, false
}

// ByOrderId searches an order by its OrderId
// returns false as second argument if the order was not found
func (s *Storage) ByOrderId(orderId string) (*entities.Order, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, order := range s.buffer {
		if order.OrderId == orderId {
			return order, true
		}
	}
	return nil, false
}
