package cartstore

import (
	"context"
	"sync"

	pb "github.com/kurtosis-tech/online-boutique-demo/cartservice/proto"
)

type memoryCartStore struct {
	sync.RWMutex

	carts map[string]map[string]*itemDetail
}

type itemDetail struct {
	quantity   int32
	isAPresent bool
}

func NewMemoryCartStore() CartStore {
	return &memoryCartStore{
		carts: make(map[string]map[string]*itemDetail),
	}
}

func (s *memoryCartStore) AddItem(_ context.Context, userID, productID string, quantity int32, isAPresent bool) error {
	s.Lock()
	defer s.Unlock()

	itemDetailObj := &itemDetail{quantity: quantity, isAPresent: isAPresent}
	if cart, ok := s.carts[userID]; ok {
		if currentItemDetail, ok := cart[productID]; ok {
			itemDetailObj.quantity = currentItemDetail.quantity + quantity
		}
		cart[productID] = itemDetailObj

		s.carts[userID] = cart
	} else {
		s.carts[userID] = map[string]*itemDetail{productID: itemDetailObj}
	}
	return nil
}

func (s *memoryCartStore) EmptyCart(_ context.Context, userID string) error {
	s.Lock()
	defer s.Unlock()

	delete(s.carts, userID)
	return nil
}

func (s *memoryCartStore) GetCart(_ context.Context, userID string) (*pb.Cart, error) {
	s.RLock()
	defer s.RUnlock()

	if cart, ok := s.carts[userID]; ok {
		items := make([]*pb.CartItem, 0, len(cart))
		for p, currentItemDetail := range cart {
			items = append(items, &pb.CartItem{ProductId: p, Quantity: currentItemDetail.quantity, IsAPresent: currentItemDetail.isAPresent})
		}
		return &pb.Cart{UserId: userID, Items: items}, nil
	}
	return &pb.Cart{UserId: userID}, nil
}
