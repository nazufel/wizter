package main

import (
	"bufio"
	"log"
	"net"
	"os"
	"strings"

	pb "github.com/nazufel/wizter/wizard"
	"google.golang.org/grpc"
)

const configMapFile = "./wizards-server-configMap.txt"

type server struct {
	pb.WizardServiceServer
}

func main() {
	log.Println("starting wizards server")

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

	log.Printf("printing MONGO_HOST: %v", os.Getenv("MONGO_HOST"))

	s, err := dbConnect()
	if err != nil {
		log.Fatalf("cannot connect to DB", err)
	}

	log.Println("dropping the wizards collection and seeding database")

	err = s.seedData()
	if err != nil {
		log.Fatalf("failed to seed the DB: %v", err)
	}
	log.Printf("finished seeding the database")

	lis, err := net.Listen("tcp", ":"+os.Getenv("WIZARDS_SERVER_GRPC_PORT"))
	if err != nil {
		log.Fatalf("could not listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	pb.RegisterWizardServiceServer(grpcServer, &server{})

	log.Println("# ----------------------------------- #")
	log.Printf("running grpc server on port: %v", os.Getenv("WIZARDS_SERVER_GRPC_PORT"))
	log.Println("# ----------------------------------- #")
	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatalf("failed to start the grpc server: %s", err)
		log.Fatalln("trivial")
	}
}

//List is the gRPC service method that retrieves a list of wizards from the database and streams back to the client
func (s *server) List(e *pb.Wizard, srv pb.WizardService_ListServer) error {

	var st storage

	wizards, err := st.getWizards()
	if err != nil {
		log.Println(err)
		return err
	}

	log.Println("# -------------------------------------- #")
	log.Println("sending list of wizards to client")
	log.Println("# -------------------------------------- #")

	for i := range wizards {
		// Commenting out. Uncomment when making the server change in the demo
		// switch house := wizard.House; house {
		// case "Gryffindor":
		// 	log.Printf("%v - sending wizard to client: %v", emoji.Eagle, wizard.GetName())
		// case "Ravenclaw":
		// 	log.Printf("%v - sending wizard to client: %v", emoji.Bird, wizard.GetName())
		// case "Hufflepuff":
		// 	log.Printf("%v - sending wizard to client: %v", emoji.Badger, wizard.GetName())
		// case "Slytherin":
		// 	log.Printf("%v - sending wizard to client: %v", emoji.Snake, wizard.GetName())
		// }
		err := srv.Send(wizards[i])
		if err != nil {
			log.Printf("send error: %v", err)
		}
	}

	// uncomment before building the docker image with Bazel
	// log.Println("")
	// log.Println("built with Bazel")
	// log.Println("")

	log.Println("# -------------------------------------- #")
	log.Println("done sending list of wizards to client")
	log.Println("# -------------------------------------- #")

	return nil

}

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
