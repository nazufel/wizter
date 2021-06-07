package main

import (
	"context"
	"log"
	"net"
	"os"
	"time"

	pb "github.com/nazufel/telepresence-demo/wizard"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"google.golang.org/grpc"
)

const (
	grpcPort = ":9999"
)

var dbConnection = "mongodb://" + os.Getenv("MONGO_HOST") + ":27017" + "/" + os.Getenv("MONGO_DATABASE")

type server struct {
	pb.UnimplementedWizardServiceServer
}

func main() {
	log.Println("starting wizards server...")
	log.Println("trying to connect to database")

	err := storageConnect()
	if err != nil {
		log.Fatalf("failed to connect to storage: %v", err)
	}

	err = seedData()
	if err != nil {
		log.Fatalf("failed to seed the DB: %v", err)
	}

	lis, err := net.Listen("tcp", grpcPort)
	if err != nil {
		log.Fatalf("could not listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	pb.RegisterWizardServiceServer(grpcServer, &server{})

	log.Printf("running grpc server on port: %v", grpcPort)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to start the grpc server: %s", err)
	}
}

var wizardCollection *mongo.Collection

func storageConnect() error {

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dbConnection))
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatalf("error pinging DB: %v", err)
	}

	wizardCollection := client.Database("wizards").Collection("wizards")

	log.Printf("connected to collection: %v", wizardCollection)

	return err
}

type WizardRecord struct {
	Name       string `bson:"name"`
	House      string `bson:"house"`
	DeathEater bool   `bson:"death_eater"`
}

func seedData() error {

	log.Println("dropping the wizards collection")

	// drop the collection in order to see fresh data for a new run
	err := wizardCollection.Drop(context.Background())
	if err != nil {
		log.Fatal("unable drop the wizard collection")
	}

	// seed the database

	wizards := []WizardRecord{
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
		{Name: "Draco Malfoy",
			House:      "Slytherin",
			DeathEater: true,
		},
		{Name: "Cho Chang",
			House:      "Raven Claw",
			DeathEater: false,
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	for i := range wizards {
		_, err = wizardCollection.InsertOne(ctx, wizards[i])
		if err != nil {
			log.Printf("failed to insert document: %v", err)
		}
		log.Printf("inserted document for wizard: %s", wizards[1].Name)
	}

	log.Printf("finished seeding the database")
	log.Printf("inserted %v documents", len(wizards))

	return err
}

// func (s *server) List(ctx context.Context, e *pb.EmptyRequest) *pb.WizardService_ListServer {

// 	var err error

// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()

// 	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dbConnection))
// 	defer func() {
// 		if err = client.Disconnect(ctx); err != nil {
// 			panic(err)
// 		}
// 	}()
// 	collection := client.Database("wizards").Collection("wizards")

// 	log.Printf("collection: %v", collection)
// 	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()

// 	err = client.Ping(ctx, readpref.Primary())
// 	if err != nil {
// 		log.Fatalf("error pinging DB: %v", err)
// 	}
// 	return nil
// }
