package main

import (
	"bufio"
	"context"
	"log"
	"net"
	"os"
	"strings"
	"time"

	pb "github.com/nazufel/wizter/wizard"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
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

	log.Println("dropping the wizards collection and seeding database")

	err = seedData()
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
func (s *server) List(e *pb.EmptyRequest, srv pb.WizardService_ListServer) error {

	log.Println("# -------------------------------------- #")
	log.Println("sending list of wizards to client")
	log.Println("# -------------------------------------- #")

	dbConnectionString := "mongodb://" + os.Getenv("MONGO_HOST") + ":" + "27017" + "/" + os.Getenv("MONGO_DATABASE")

	client, err := mongo.NewClient(options.Client().ApplyURI(dbConnectionString))

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	defer func() {
		err = client.Disconnect(ctx)
		if err != nil {
			log.Fatalf("failed to disconnect client:  %v", err)
		}
	}()

	// test connection
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatalf("error pinging DB: %v", err)
	}

	db := client.Database("wizards")

	// set find options behavior
	findOptions := options.Find()
	// findOptions.SetLimit(25)

	// filter by company, this is the default behavior
	filter := bson.D{{}}

	cursor, err := db.Collection("wizards").Find(ctx, filter, findOptions)
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

// seedData drops the wizards collection and seeds it with fresh data to the demo
func seedData() error {

	dbConnectionString := "mongodb://" + os.Getenv("MONGO_HOST") + ":" + "27017" + "/" + os.Getenv("MONGO_DATABASE")

	client, err := mongo.NewClient(options.Client().ApplyURI(dbConnectionString))

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	defer func() {
		err = client.Disconnect(ctx)
		if err != nil {
			log.Fatalf("failed to disconnect client:  %v", err)
		}
	}()

	// test connection
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatalf("error pinging DB: %v", err)
	}

	db := client.Database("wizards")

	// // drop the collection in order to see fresh data for a new run
	err = db.Collection("wizards").Drop(context.Background())
	if err != nil {
		log.Fatal("unable drop the wizard collection")
	}

	// seed the database

	wizards := []pb.Wizard{
		{Name: "Harry Potter",
			House:      "Gryffindor",
			DeathEater: false,
		},
		{Name: "Ron Weasley",
			House:      "Gryffindor",
			DeathEater: false,
		},
		{Name: "Hermione Granger",
			House:      "Gryffindor",
			DeathEater: false,
		},
		{Name: "Cho Chang",
			House:      "Ravenclaw",
			DeathEater: false,
		},
		{Name: "Luna Lovegood",
			House:      "Ravenclaw",
			DeathEater: false,
		},
		{Name: "Sybill Trelawney",
			House:      "Ravenclaw",
			DeathEater: false,
		},
		{Name: "Pomona Sprout",
			House:      "Hufflepuff",
			DeathEater: false,
		},
		{Name: "Cedric Diggory",
			House:      "Hufflepuff",
			DeathEater: false,
		},
		{Name: "Newton Scamander",
			House:      "Hufflepuff",
			DeathEater: false,
		},
		{Name: "Draco Malfoy",
			House:      "Slytherin",
			DeathEater: true,
		},
		{Name: "Bellatrix Lestrange",
			House:      "Slytherin",
			DeathEater: true,
		},
		{Name: "Severus Snape",
			House:      "Slytherin",
			DeathEater: false,
		},
	}
	log.Printf("connected to the database")

	for i := range wizards {
		_, err = db.Collection("wizards").InsertOne(ctx, wizards[i])
		if err != nil {
			log.Fatalf("failed to insert document: %v", err)
		}
		log.Printf("inserted document for wizard: %s", wizards[i].Name)
		i++
	}

	log.Printf("inserted %v documents", len(wizards))

	return err
}
