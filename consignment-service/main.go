package main

import (
	pb "github.com/ericnjoroge/shippy-microservices/consignment-service"
)

const (
	port = ":50051"
)

type repository interface {
	Create(*pb.Consignment) (*pb.Consignment, error)
}
