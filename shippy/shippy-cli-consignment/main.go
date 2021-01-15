package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"context"

	pb "github.com/samueldsr/shippy-service-consignment/proto/consignment"
	"google.golang.org/grpc"
)

const (
	address         = "localhost:50051"
	defaultFilename = "consignment.json"
)

func parseFile(file string) (*pb.Consignement, error) {
	var consignment *pb.Consignement
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(data, &consignment)
	return consignment, err
}

func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Did not connect: %v", err)
	}

	defer conn.Close()

	client := pb.NewShippingServiceClient(conn)

	file := defaultFilename

	if len(os.Args) > 1 {
		file = os.Args[1]
	}

	consignment, err := parseFile(file)
	if err != nil {
		log.Fatalf("Could not open file: %v", file)
	}

	r, err := client.CreateConsignement(context.Background(), consignment)
	if err != nil {
		log.Fatalf("Could not greet: %v", err)
	}
	log.Printf("Created: %t", r.Created)

	all, err := client.GetConsignements(context.Background(), &pb.GetRequest{})
	if err != nil {
		log.Fatalf("Couldn't list all consignments: %v: %T", err, err)
	}
	for _, v := range all.Consignements {
		log.Printf("%T: %v", v, v.GetId())
	}

}
