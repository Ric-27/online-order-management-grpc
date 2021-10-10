package product

import (
	"context"

	productrpc "github.com/Ric-27/online-order-management-grpc/product/rpc"

	"github.com/Ric-27/online-order-management-grpc/store"
)

// Service holds RPC handlers for the product service. It implements the product.ServiceServer interface.
type service struct {
	productrpc.UnimplementedServiceServer
	s store.ProductStore
}

func NewService(s store.ProductStore) *service {
	return &service{s: s}
}

// Fetch all existing products in the system.
func (s *service) ListProducts(ctx context.Context, r *productrpc.ListProductsRequest) (*productrpc.ListProductsResponse, error) {
	return &productrpc.ListProductsResponse{Products: s.s.Products()}, nil
}
func (s *service) ProductOfId(ctx context.Context, r *productrpc.ProductOfIdRequest) (*productrpc.ProductOfIdResponse, error) {
	prod, err := s.s.Product(r.Id)
	return &productrpc.ProductOfIdResponse{Product: prod}, err
}
