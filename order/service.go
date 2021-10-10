package order

import (
	"context"

	orderrpc "github.com/Ric-27/online-order-management-grpc/order/rpc"
	"github.com/Ric-27/online-order-management-grpc/store"
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
