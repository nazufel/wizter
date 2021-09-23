package main

import (
	"context"
	"io"
	"log"
	"os"
	"time"

	pb "github.com/nazufel/wizter/wizard"

	"google.golang.org/grpc"
)

func main() {

	loadConfigs()

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

	address := os.Getenv("WIZARDS_SERVER_SERVICE_HOST") + ":" + os.Getenv("WIZARDS_SERVER_GRPC_PORT")

	log.Println("starting client")

	log.Printf("connecting to address: %v", address)
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

		// commenting out for the demo. uncomment during demo of the client
		// if resp.GetDeathEater() {
		// 	alertDeathEather(resp)
		// }
	}
}

// commenting out for the demo. uncomment during demo of the client
// func alertDeathEather(w *pb.Wizard) {
// 	log.Println("")
// 	log.Printf("Oh no! %s is a Death Eater!", w.GetName())
// 	log.Println("")
// }

// loadConfigs looks for a specific file of key=value pairs and loads them as variables for the runtime instance
func loadConfigs() {

	if os.Getenv("WIZARDS_SERVER_SERVICE_HOST") == "" {
		os.Setenv("WIZARDS_SERVER_SERVICE_HOST", "localhost")
	}

	if os.Getenv("WIZARDS_SERVER_GRPC_PORT") == "" {
		os.Setenv("WIZARDS_SERVER_GRPC_PORT", "9999")
	}
}