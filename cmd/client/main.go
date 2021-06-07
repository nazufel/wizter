package main

import (
	"context"
	"log"
	"time"

	pb "github.com/nazufel/telepresence-demo/wizard"

	"google.golang.org/grpc"
)

const (
	address = "localhost:9999"
)

func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewWizardServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)

	defer cancel()

	r, err := c.Add(ctx, &pb.Wizard{
		Name:       "Harry Potter",
		House:      "Gryffindor",
		DeathEater: true,
	})
	if err != nil {
		log.Fatalf("could not send message: %v", err)
	}

	wizard := pb.Wizard{
		Id:         r.GetId(),
		Name:       r.GetName(),
		House:      r.GetHouse(),
		DeathEater: r.GetDeathEater(),
	}

	log.Printf("received from server: %s", &wizard)
}
