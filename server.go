package main

import (
	"log"
	"net"
	"os"
	"os/signal"

	"github.com/Ric-27/online-order-management-grpc/order"
	orderpb "github.com/Ric-27/online-order-management-grpc/order/rpc"
	"github.com/Ric-27/online-order-management-grpc/product"
	productpb "github.com/Ric-27/online-order-management-grpc/product/rpc"
	"github.com/Ric-27/online-order-management-grpc/store"
	"google.golang.org/grpc"
)

func main() {
	log := log.Default()
	ps := store.NewProductStore()
	ords := store.NewOrderStore()

	lis, err := net.Listen("tcp", "0.0.0.0:8000")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	defer lis.Close()

	opts := []grpc.ServerOption{}
	srv := grpc.NewServer(opts...)
	defer srv.Stop()
	//declaration of services
	prd := product.NewService(ps)
	productpb.RegisterServiceServer(srv, prd)

	ord := order.NewService(ords, ps)
	orderpb.RegisterServiceServer(srv, ord)

	go func() {
		log.Println("server started on port 8000...")
		if err := srv.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	// wait for control C to exit
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	<-ch
	log.Println("stopping the server")
}
