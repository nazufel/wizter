package main

import (
	"context"
	"io"
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

	stream, err := c.List(ctx, &pb.EmptyRequest{})
	if err != nil {
		log.Fatalf("could not send message: %v", err)
	}

	done := make(chan bool)

	go func() {
		for {
			resp, err := stream.Recv()
			if err == io.EOF {
				done <- true
				return
			}
			if err != nil {
				log.Fatalf("cannot receive stream: %v", err)
			}
			log.Printf("Wizard received: %v", resp.GetName())
		}
	}()

	<-done
	log.Printf("client finished")
}
