package main

import (
	"context"
	"fmt"
	"log"

	micro "github.com/micro/go-micro"
	pb "github.com/vodaza36/go-micro/consignment-service/proto/consignment"
	vesselProto "github.com/vodaza36/go-micro/vessel-service/proto/vessel"
)

const (
	port = ":50051"
)

// IRepository repo interface
type IRepository interface {
	Create(*pb.Consignment) (*pb.Consignment, error)
	GetAll() []*pb.Consignment
}

// Repository - Dummy repository, this simulates the use of a datastore
// of some kind. We'll replace this with a real implementation later on.
type Repository struct {
	consignments []*pb.Consignment
}

// Create a new consignment record in the repository
func (r *Repository) Create(consignment *pb.Consignment) (*pb.Consignment, error) {
	updated := append(r.consignments, consignment)
	r.consignments = updated
	return consignment, nil
}

// GetAll returns all consignments
func (r *Repository) GetAll() []*pb.Consignment {
	return r.consignments
}

// Service should implement all of the methods to satisfy the service
// we defined in our protobuf definition. You can check the interface
// in the generated code itself for the exact method signatures etc
// to give you a better idea.
type service struct {
	repo         IRepository
	vesselClient vesselProto.VesselServiceClient
}

// CreateConsignment - we created just one method on our service,
// which is a create method, which takes a context and a request as an
// argument, these are handled by the gRPC server.
// CreateConsignment - we created just one method on our service,
// which is a create method, which takes a context and a request as an
// argument, these are handled by the gRPC server.
func (s *service) CreateConsignment(ctx context.Context, req *pb.Consignment, res *pb.Response) error {
	// Here we call a client instance of our vessel service with our consignment weight,
	// and the amount of containers as the capacity value
	vesselResponse, err := s.vesselClient.FindAvailable(context.Background(), &vesselProto.Specification{
		MaxWeight: req.Weight,
		Capacity:  int32(len(req.Containers)),
	})
	log.Printf("Found vessel: %s \n", vesselResponse.Vessel.Name)
	if err != nil {
		return err
	}
	// We set the VesselId as the vessel we got back from our
	// vessel service
	req.VesselId = vesselResponse.Vessel.Id
	// Save our consignment
	consignment, err := s.repo.Create(req)
	if err != nil {
		return err
	}
	// Return matching the `Response` message we created in our
	// protobuf definition.
	res.Created = true
	res.Consignment = consignment
	return nil
}

// GetConsignments retuns an array of consignments
func (s *service) GetConsignments(ctx context.Context, req *pb.GetRequest, res *pb.Response) error {
	consignments := s.repo.GetAll()
	res.Consignments = consignments
	return nil
}

func main() {
	repo := &Repository{}
	// Create a new service. Optionally include some options here.
	srv := micro.NewService(
		// This name must match the package name given in your protobuf definition
		micro.Name("go.micro.srv.consignment"),
		micro.Version("latest"),
	)
	vesselClient := vesselProto.NewVesselServiceClient("go.micro.srv.vessel", srv.Client())
	// Init will parse the command line flags.
	srv.Init()
	// Register handler
	pb.RegisterShippingServiceHandler(srv.Server(), &service{repo, vesselClient})
	// Run the server
	if err := srv.Run(); err != nil {
		fmt.Println(err)
	}
}
