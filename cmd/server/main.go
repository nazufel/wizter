package main

import (
	"context"
	"log"
	"net"
	"os"
	"time"

	pb "github.com/nazufel/telepresence-demo/wizard"
	"go.mongodb.org/mongo-driver/bson"
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
	pb.WizardServiceServer
}

func main() {
	log.Println("starting wizards server")

	log.Println("dropping the wizards collection and seeding database")
	err := seedData()
	if err != nil {
		log.Fatalf("failed to seed the DB: %v", err)
	}
	log.Printf("finished seeding the database")

	lis, err := net.Listen("tcp", grpcPort)
	if err != nil {
		log.Fatalf("could not listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	pb.RegisterWizardServiceServer(grpcServer, &server{})

	log.Printf("running grpc server on port: %v", grpcPort)
	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatalf("failed to start the grpc server: %s", err)
	}
}

type WizardRecord struct {
	Name       string `bson:"name"`
	House      string `bson:"house"`
	DeathEater bool   `bson:"death_eater"`
}

func (s *server) List(e *pb.EmptyRequest, srv pb.WizardService_ListServer) error {

	log.Println("sending list of wizards to client")

	client, err := mongo.NewClient(options.Client().ApplyURI(dbConnection))

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

	cursor, err := db.Collection("Wizards").Find(ctx, bson.M{})
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

		err = srv.Send(&wizard)
		if err != nil {
			log.Printf("send error: %v", err)
		}
		log.Printf("sent wizard to client: %v", wizard.GetName())
	}

	err = cursor.Err()
	if err != nil {
		log.Printf("error with the client cursor: %v", err)
	}

	return nil

}

// seedData drops the wizards collection and seeds it with fresh data to the demo
func seedData() error {

	client, err := mongo.NewClient(options.Client().ApplyURI(dbConnection))

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
