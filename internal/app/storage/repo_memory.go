package storage

import (
	"sync"
)

type OrderRepoMemory struct {
	mutex        sync.RWMutex
	idToOrderMap map[string]Order
}

func NewOderRepoMemory() *OrderRepoMemory {
	return &OrderRepoMemory{
		mutex:        sync.RWMutex{},
		idToOrderMap: make(map[string]Order),
	}
}

func (r *OrderRepoMemory) SaveOrder(order Order) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, ok := r.idToOrderMap[order.OrderUID]; ok {
		return ErrOrderExists
	}
	r.idToOrderMap[order.OrderUID] = order
	return nil
}

func (r *OrderRepoMemory) GetOrderByID(orderID string) (Order, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	if _, ok := r.idToOrderMap[orderID]; !ok {
		return Order{}, ErrOrderNotExists
	}

	return r.idToOrderMap[orderID], nil
}
