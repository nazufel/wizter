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
	pb.UnimplementedWizardServiceServer
}

func main() {
	log.Println("starting wizards server...")
	log.Println("trying to connect to database")

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

type WizardRecord struct {
	Name       string `bson:"name"`
	House      string `bson:"house"`
	DeathEater bool   `bson:"death_eater"`
}

func (s *server) Add(ctx context.Context, wz *pb.Wizard) (*pb.Wizard, error) {

	log.Printf("receved wizard: %s", wz)

	var err error

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dbConnection))
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
	collection := client.Database("wizards").Collection("wizards")

	log.Printf("collection: %v", collection)
	ctx, cancel = context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatalf("error pinging DB: %v", err)
	}

	_, err = collection.InsertOne(ctx, wz)
	if err != nil {
		log.Printf("failed to insert document: %v", err)
	}

	wizardrecord := WizardRecord{}

	filter := bson.D{{Key: "name", Value: wz.GetName()}}
	ctx, cancel = context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	err = collection.FindOne(ctx, filter).Decode(&wizardrecord)

	response := &pb.Wizard{
		Name:       wizardrecord.Name,
		House:      wizardrecord.House,
		DeathEater: wizardrecord.DeathEater,
	}

	return response, err
}

func (s *server) List(ctx context.Context, wz *pb.Wizard) (*pb.Wizard, error) {

	var err error

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dbConnection))
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
	collection := client.Database("wizards").Collection("wizards")

	log.Printf("collection: %v", collection)
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatalf("error pinging DB: %v", err)
	}
	return wz, nil
}
