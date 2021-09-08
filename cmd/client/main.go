package main

import (
	"bufio"
	"context"
	"io"
	"log"
	"os"
	"strings"
	"time"

	pb "github.com/nazufel/wizter/wizard"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const configMapFile = "./wizards-server-configMap.txt"

// clientCall func makes the call to the sever to get a list of wizards
func clientCall(ic int32) {

	address := os.Getenv("WIZARDS_SERVER_SERVICE_HOST") + ":" + os.Getenv("WIZARDS_SERVER_GRPC_PORT")

	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewWizardServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// adding Telepresence intercept header to set up a personal intercept - https://www.getambassador.io/docs/telepresence/latest/howtos/intercepts/#4-personal-intercept
	// the client will only progate headers on everyother request so that the user will only be able to intercept up to half of the requests and this won't cause a
	// complete outage.
	if ic%2 == 0 {
		md := metadata.New(map[string]string{"x-telepresence-intercept-id": "tele-demo"})
		ctx = metadata.NewOutgoingContext(ctx, md)
	}

	stream, err := c.List(ctx, &pb.Wizard{})
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

	file, err := os.Open(configMapFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		// fmt.Println(scanner.Text())
		s := strings.Split(scanner.Text(), "=")
		os.Setenv(s[0], s[1])
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	log.Println("done setting environment variables")

}

// func main is the main function
func main() {
	// set envs for local dev if the intercept env file exists
	log.Printf("checking for configMap file at: %v", configMapFile)
	_, err := os.Stat(configMapFile)
	if !os.IsNotExist(err) {
		log.Printf("found %s file. setting environment variables", configMapFile)
		loadConfigs()
	}
	if os.IsNotExist(err) {
		log.Printf("did not find config file: %s. using Kubernetes environment", configMapFile)
	}

	// interceptCounter increases a counter to control if
	var interceptCounter int32

	for {
		time.Sleep(1 * time.Second)
		log.Println("# -------------------------------------- #")
		log.Println("requesting list of wizards from the server")
		log.Println("# -------------------------------------- #")

		// call for a list of wizards and pass in the counter
		clientCall(interceptCounter)
		log.Println("# --------------------------------------#")
		log.Println("received list of wizards from the server")
		log.Println("# ------------------------------------- #")
		// if counter reaches greather than 100,000,000, then reset it as to not overflow
		if interceptCounter > 100000000 {
			interceptCounter = 0
		}
		interceptCounter++
	}
}
