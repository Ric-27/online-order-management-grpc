package order

import (
	"context"

	orderrpc "github.com/bigbluedisco/tech-challenge/backend/v1/order/rpc"
	"github.com/bigbluedisco/tech-challenge/backend/v1/store"
)

// Service holds RPC handlers for the order service. It implements the orderrpc.ServiceServer interface.
type service struct {
	orderrpc.UnimplementedServiceServer
	s store.OrderStore
}

func NewService(s store.OrderStore) *service {
	return &service{s: s}
}

// Fetch all existing orders in the system.
func (s *service) ListOrders(ctx context.Context, r *orderrpc.ListOrdersRequest) (*orderrpc.ListOrdersResponse, error) {
	return &orderrpc.ListOrdersResponse{Orders: s.s.Orders()}, nil
}
