package order

import (
	"encoding/json"
	"errors"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"context"

	orderrpc "github.com/Ric-27/online-order-management-grpc/order/rpc"

	"github.com/Ric-27/online-order-management-grpc/store"
)

// Service holds RPC handlers for the order service. It implements the orderrpc.ServiceServer interface.
type service struct {
	orderrpc.UnimplementedServiceServer
	s  store.OrderStore
	s2 store.ProductStore
}

func NewService(s store.OrderStore, s2 store.ProductStore) *service {
	return &service{s: s, s2: s2}
}

// Fetch all existing orders in the system.
func (s *service) ListOrders(ctx context.Context, r *orderrpc.ListOrdersRequest) (*orderrpc.ListOrdersResponse, error) {
	return &orderrpc.ListOrdersResponse{Orders: s.s.Orders()}, nil
}

//fields of importance and their sub groups division of the address api response
type Properties struct {
	Name     string
	Postcode string
	City     string
	Score    float64
}

type Features struct {
	Properties Properties
}

type TargetResponse struct {
	Features []Features
}

//perform an http get request and transform the answer onto a json (from )
// from https://stackoverflow.com/questions/17156371/how-to-get-json-response-from-http-get
var myClient = &http.Client{Timeout: 10 * time.Second}

func getJson(url string, target interface{}) error {
	r, err := myClient.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(target)
}

func gen_id(lName string) (id string) {
	//create an unique identifier for the order, composed of a random 4 digit number, the date and time, and the person last name
	id = strconv.Itoa(rand.Intn(9000)+1000) + "-" + time.Now().Format("20060102-150405") + "-" + strings.ToUpper(lName)
	return
}

const APIADDPREFIX string = "https://api-adresse.data.gouv.fr/search/?q="
const APIADDSUFIX string = "&limit=1"
const MIN_SCORE float32 = 0.3

func (s *service) CreateOrder(ctx context.Context, r *orderrpc.CreateOrderRequest) (*orderrpc.CreateOrderResponse, error) {
	//check if the products are valid, if not return an error
	for key := range r.ProductList {
		_, err := s.s2.Product(key)
		if err != nil {
			return &orderrpc.CreateOrderResponse{Order: nil}, errors.New("product " + key + " not available in store, order was not created")
		}
	}
	//consult on the address api if the address is correct
	resp := TargetResponse{} //expected body answer from the api
	linename := strings.Replace(r.Address.Street, " ", "+", -1)
	city := strings.Replace(r.Address.City, " ", "+", -1)
	txt_request := linename + ",+" + r.Address.Postal + ",+" + city
	request := APIADDPREFIX + txt_request + APIADDSUFIX
	getJson(request, &resp) //request to the api and transform the answer into a map
	//if the api returns an error or the address has a score less than 0.3 return an error
	if len(resp.Features) == 0 || resp.Features[0].Properties.Score < float64(MIN_SCORE) {
		return &orderrpc.CreateOrderResponse{Order: nil}, errors.New("address not found")
	}
	address := orderrpc.Address{Street: resp.Features[0].Properties.Name, Postal: resp.Features[0].Properties.Postcode, City: resp.Features[0].Properties.City, Country: "France"}
	//create an unique id
	gid := gen_id(r.LName)
	//order object created
	ord := orderrpc.Order{
		Id:          gid,
		FName:       r.FName,
		LName:       r.LName,
		ProductList: r.ProductList,
		Address:     &address,
	}
	s.s.SetOrder(&ord) //order added to the system
	return &orderrpc.CreateOrderResponse{Order: &ord}, nil
}
