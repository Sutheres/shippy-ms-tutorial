package main

import (
	pb "github.com/sutheres/shippy-service-consignment/proto/consignment"
)

const (
	port = ":50051"
)

type repository interface {
	Create(*pb.Consignment)
}