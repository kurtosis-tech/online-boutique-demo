package cartstore

import (
	"context"

	pb "github.com/kurtosis-tech/online-boutique-demo/cartservice/proto"
)

type CartStore interface {
	AddItem(ctx context.Context, userID, productID string, quantity int32) error
	EmptyCart(ctx context.Context, userID string) error
	GetCart(ctx context.Context, userID string) (*pb.Cart, error)
}
