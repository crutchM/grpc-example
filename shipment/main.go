package main

import (
	"awesomeProject/shipment/rpc"
	context "context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

const (
	port = ":50051"
)

type IRepository interface {
	Create(consignment *rpc.Consignment) (*rpc.Consignment, error)
	GetAll() []*rpc.Consignment
}

type Repository struct {
	consignments []*rpc.Consignment
}

func (s *Repository) Create(consignment *rpc.Consignment) (*rpc.Consignment, error) {
	updated := append(s.consignments, consignment)
	s.consignments = updated
	return consignment, nil
}

func (s *Repository) GetAll() []*rpc.Consignment {
	return s.consignments
}

type ShipmentService struct {
	repo IRepository
}

func (s ShipmentService) CreateConsignment(ctx context.Context, consignment *rpc.Consignment) (*rpc.Response, error) {
	consignment, err := s.repo.Create(consignment)
	if err != nil {
		return nil, err
	}

	return &rpc.Response{Created: true, Consignment: consignment}, nil
}

func (s ShipmentService) GetConsignments(ctx context.Context, request *rpc.GetRequest) (*rpc.Response, error) {
	consignments := s.repo.GetAll()
	return &rpc.Response{Consignments: consignments}, nil
}

func main() {
	repo := &Repository{}

	lis, err := net.Listen("tcp", port)

	if err != nil {
		log.Fatal(err)
	}

	s := grpc.NewServer()

	rpc.RegisterShippingServiceServer(s, &ShipmentService{repo})

	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
