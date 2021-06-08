package main

import (
	"context"
	"io"
	"log"
	"os"
	"time"

	pb "github.com/nazufel/telepresence-demo/wizard"

	"google.golang.org/grpc"
)

var address = os.Getenv("SERVER_HOST") + ":" + os.Getenv("GRPC_PORT")

func main() {

	log.Printf("connection address: %s", address)

	for {
		time.Sleep(1 * time.Second)
		log.Println("# -------------------------------------- #")
		log.Println("requesting list of wizards from the server")
		log.Println("# -------------------------------------- #")

		clientCall()
		log.Println("# --------------------------------------#")
		log.Println("received list of wizards from the server")
		log.Println("# ------------------------------------- #")
	}
}

func clientCall() {

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

	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			return
		}
		if err != nil {
			log.Fatalf("cannot receive stream: %v", err)
		}
		log.Printf("wizard received: %v", resp.GetName())
	}
}
