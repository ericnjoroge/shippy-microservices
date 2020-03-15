package main

import (
	"context"
	"log"
	"net"
	"sync"

	pb "github.com/ericnjoroge/shippy-microservices/consignment-service/proto/consignment"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	port = ":50051"
)

type repository interface {
	Create(*pb.Consignment) (*pb.Consignment, error)
	GetAll() []*pb.Consignment
}

//Repository - Dummy repository that simulates the use of a datastore
type Repository struct {
	mu           sync.RWMutex
	consignments []*pb.Consignment
}

//Create a new consignment
func (repo *Repository) Create(consignment *pb.Consignment) (*pb.Consignment, error) {
	repo.mu.Lock()
	updated := append(repo.consignments, consignment)
	repo.consignments = updated
	repo.mu.Unlock()
	return consignment, nil
}

// Get all consignments
func (repo *Repository) GetAll() []*pb.Consignment {
	return repo.consignments
}

// Service implements all of the methods to satisfy the service defined in the protobuf definition
type service struct {
	repo repository
}

//CreateConsignment method - creates a single method on the service
// create method takes a context and a request as an argument that are handled by the grpc server
func (s *service) CreateConsignment(ctx context.Context, req *pb.Consignment) (*pb.Response, error) {

	//save the consignment
	consignment, err := s.repo.Create(req)
	if err != nil {
		return nil, err
	}

	//return a response matching the `Response` message created in the protobuf definition
	return &pb.Response{Created: true, Consignment: consignment}, nil
}

// Get all consignments method to fetch all created consignments
func (s *service) GetConsignments(ctx context.Context, req *pb.GetRequest) (*pb.Response, error) {
	consignments := s.repo.GetAll()
	return &pb.Response{Consignments: consignments}, nil
}

func main() {

	repo := &Repository{}

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("[main func] failed to listen: %v", err)
	}
	s := grpc.NewServer()

	//Register the service with the gRPC server, to tie the auto-generated interface code for the protobuf definition
	pb.RegisterShippingServiceServer(s, &service{repo})

	//register the reflection service on the gRPC server
	reflection.Register(s)

	log.Println("Running on port:", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
