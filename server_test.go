package main

import (
	"context"
	"log"
	"net"
	"testing"

	"github.com/Ric-27/online-order-management-grpc/order"
	orderpb "github.com/Ric-27/online-order-management-grpc/order/rpc"
	"github.com/Ric-27/online-order-management-grpc/product"
	productpb "github.com/Ric-27/online-order-management-grpc/product/rpc"
	"github.com/Ric-27/online-order-management-grpc/store"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

// example from https://stackoverflow.com/questions/42102496/testing-a-grpc-service
//the module bufconn allows for testing a gRPC server without openit to a port
const bufSize = 1024 * 1024

var lis *bufconn.Listener

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

//instead of main we just declare the server needed for testing
func init() {
	log := log.Default()
	ps := store.NewProductStore()
	ords := store.NewOrderStore()

	lis = bufconn.Listen(bufSize)

	srv := grpc.NewServer()

	prd := product.NewService(ps)
	productpb.RegisterServiceServer(srv, prd)

	ord := order.NewService(ords, ps)
	orderpb.RegisterServiceServer(srv, ord)

	go func() {
		if err := srv.Serve(lis); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()
}

//test the service of getting a product from its id, the comparison is between the error and if the function should give out an error
func TestProductId(t *testing.T) {
	var id_test = map[string]bool{"PIPR-CRMSOL-50ML": false, "INEXISTANT-PRODUCT": true, "PIPR-JACKET-SIZM": false, "WHO-AM-I": true}
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()
	client := productpb.NewServiceClient(conn)
	for key, value := range id_test {
		_, err := client.ProductOfId(ctx, &productpb.ProductOfIdRequest{Id: key})
		if value == (err == nil) {
			t.Errorf("FAIL: %v", key)
			t.Logf("Expected: %+v | got: %+v", value, !(err == nil))
		}
	}
}

// Test the service of creating an order, the order should be refused if a chosen product is not on the store,
//or if the address cant be validated, the comparison is between the error and the fact that the function should give out and error or not
func TestOrderCreate(t *testing.T) {
	id_test := make(map[*orderpb.Order]bool)
	address1 := &orderpb.Address{Street: "8 Route Erica Tabarly", Postal: "91300", City: "Masi", Country: "France"}
	productlist1 := map[string]int32{"PIPR-CRMSOL-50ML": 3}
	order1 := orderpb.Order{FName: "Ricardo", LName: "Rico", Address: address1, ProductList: productlist1}
	id_test[&order1] = false

	address2 := &orderpb.Address{Street: "8 Route Eric Tabarly", Postal: "91300", City: "Massy", Country: "France"}
	productlist2 := map[string]int32{"INEXISTANT-PRODUCT": 3}
	order2 := orderpb.Order{FName: "prod-error", LName: "Rico", Address: address2, ProductList: productlist2}
	id_test[&order2] = true

	address3 := &orderpb.Address{Street: "29 Carrera 38", Postal: "120", City: "Kansas", Country: ""}
	productlist3 := map[string]int32{"PIPR-CRMSOL-50ML": 3}
	order3 := orderpb.Order{FName: "address-error", LName: "Rico", Address: address3, ProductList: productlist3}
	id_test[&order3] = true

	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()
	client := orderpb.NewServiceClient(conn)
	for key, value := range id_test {
		res, err := client.CreateOrder(ctx, &orderpb.CreateOrderRequest{FName: key.FName, LName: key.LName, Address: key.Address, ProductList: key.ProductList})
		if value == (err == nil) {
			t.Errorf("FAIL: %v", res)
			t.Logf("Expected: %+v | got: %+v", value, !(err == nil))
		}
	}
}
