package main

import (
	"context"
	"github.com/micro/go-micro/v2"
	pb "github.com/sutheres/shippy-ms-tutorial/shippy-service-consignment/proto/consignment"
	"log"
)

const (
	port = ":50051"
)

type repository interface {
	Create(*pb.Consignment) (*pb.Consignment, error)
	GetAll() []*pb.Consignment
}

// Repository - Dummy repository, this simulates the use of a datastore
// of some kind. We'll replace this with a real implementation later on.
type Repository struct {
	consignments []*pb.Consignment
}


func (r *Repository) Create(consignment *pb.Consignment) (*pb.Consignment, error) {
	updated := append(r.consignments, consignment)
	r.consignments = updated
	return consignment, nil
}


func (r *Repository) GetAll() []*pb.Consignment {
	return r.consignments
}

// Service should implement all of the methods to satisfy the service
// we defined in our protobuf definition. You can check the interface
// in the generated code itself for the exact method signatures etc
// to give you a better idea.
type consignmentService struct {
	repo repository
}

// CreateConsignment - we created just one method on our service,
// which is a create method, which takes a context and a request as an
// argument, these are handled by the gRPC server.
func (s *consignmentService) CreateConsignment(ctx context.Context, req *pb.Consignment, resp *pb.Response) error {

	// Save our consignment
	consignment, err := s.repo.Create(req)
	if err != nil {
		return err
	}

	// Return matching the `Response` message we created in our
	// protobuf definition
	resp.Created = true
	resp.Consignment = consignment
	return nil
}

func (s *consignmentService) GetConsignments(ctx context.Context, req *pb.GetRequest, resp *pb.Response) error {
	consignments := s.repo.GetAll()
	resp.Consignments = consignments
	return nil
}

func main() {

	repo := &Repository{}

	// Create a new service. Optionally include some options here.
	service := micro.NewService(
		micro.Name("shippy.service.consignment"),
	)

	// Init will parse the command line flags.
	service.Init()

	// Register our service with the gRPC server
	if err := pb.RegisterShippingServiceHandler(service.Server(), &consignmentService{repo}); err != nil {
		log.Panic(err)
	}

	log.Println("running on port:", port)

	if err := service.Run(); err != nil {
		log.Fatalf("failed to run: %v", err)
	}
}