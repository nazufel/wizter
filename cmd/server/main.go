package main

import (
	"context"
	"log"
	"net"
	"os"

	pb "github.com/nazufel/wizter/wizard"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
)

const configMapFile = "./wizards-server-configMap.txt"

type server struct {
	pb.WizardServiceServer
}

func main() {
	log.Println("starting wizards server")

	loadConfigs()

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

	log.Println("# -------------------------------------- #")
	log.Println("sending list of wizards to client")
	log.Println("# -------------------------------------- #")


	var st storage

	// set find options behavior
	findOptions := options.Find()
	// findOptions.SetLimit(25)

	// filter by company, this is the default behavior
	filter := bson.D{{}}

	cursor, err := st.db.Collection("wizards").Find(context.Background(), filter, findOptions)
	if err != nil {
		log.Printf("failed to get wizards: %v", err)
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var wizard pb.Wizard

		err := cursor.Decode(&wizard)
		if err != nil {
			log.Printf("unable to decode wizard cursor into struct: %v", err)
		}

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

		// comment this log statement as part of the server demo
		log.Printf("sending wizard to client: %v", wizard.GetName())

		err = srv.Send(&wizard)
		if err != nil {
			log.Printf("send error: %v", err)
		}
	}

	err = cursor.Err()
	if err != nil {
		log.Printf("error with the client cursor: %v", err)
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

// loadConfigs looks for env variables and if they aren't set, then it sets happy defaults
func loadConfigs() {

	// check and set the port for the GRPC server to listen on
	if os.Getenv("WIZARDS_SERVER_GRPC_PORT") == "" {
		os.Setenv("WIZARDS_SERVER_GRPC_PORT", "9999")
	}

	// check and set the dna name of mongo
	if os.Getenv("MONGO_HOST") == "" {
		os.Setenv("MONGO_HOST", "mongodb://mongo.default.svc.cluster.local")
	}

}
