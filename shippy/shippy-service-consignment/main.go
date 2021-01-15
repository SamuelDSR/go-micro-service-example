package main

import (
	"context"
	"log"
	"net"
	"sync"

	pb "github.com/samueldsr/shippy-service-consignment/proto/consignment"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	port = ":50051"
)

type repository interface {
	Create(*pb.Consignement) (*pb.Consignement, error)
	GetAll() []*pb.Consignement
}

type Repository struct {
	mu           sync.RWMutex
	consignments []*pb.Consignement
}

// create a new consignment
func (repo *Repository) Create(consignment *pb.Consignement) (*pb.Consignement, error) {
	repo.mu.Lock()
	updated := append(repo.consignments, consignment)
	repo.consignments = updated
	repo.mu.Unlock()
	return consignment, nil
}

// Get all consignments
func (repo *Repository) GetAll() []*pb.Consignement {
	return repo.consignments
}

type service struct {
	repo Repository
}

func (s *service) CreateConsignement(ctx context.Context,
	rep *pb.Consignement) (*pb.Response, error) {
	// Save our passed consignment
	consignment, err := s.repo.Create(rep)
	if err != nil {
		return nil, err
	}

	// Return matching the `Response` message we created in our pb definition
	return &pb.Response{Created: true, Consignement: consignment}, nil
}

func (s *service) GetConsignements(ctx context.Context,
	req *pb.GetRequest) (*pb.Response, error) {
	consignements := s.repo.GetAll()
	return &pb.Response{Consignements: consignements}, nil
}

func main() {
	repo := Repository{}

	lis, err := net.Listen("tcp", port)

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()

	pb.RegisterShippingServiceServer(s, &service{repo})

	reflection.Register(s)

	log.Println("Running on port: ", port)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to server: %v", err)
	}

}
