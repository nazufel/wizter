package main

import (
	"context"
	"log"
	"net"
	"os"
	"time"

	pb "github.com/nazufel/telepresence-demo/pkg/proto/wizard"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
)

const (
	grpcPort = ":9999"
)

var dbConnection = "mongodb://" + os.Getenv("MONGO_USER") + ":" + os.Getenv("MONGO_PASSWORD") + "@" + os.Getenv("MONGO_HOST") + ":27017" + "/" + os.Getenv("MONGO_DATABASE") + "?authSource=admin"

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

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err = collection.InsertOne(ctx, wz)
	if err != nil {
		log.Printf("failed to insert document: %v", err)
	}
	return wz, err
}

func (s *server) List(ctx context.Context, wz *pb.Wizard) (*pb.Wizard, error) {
	if wz.GetDeathEater() {
		log.Printf("Wizard: %s is a DeathEater! Run!", wz.GetName())
		return wz, nil
	}
	log.Printf("%s is not a DeathEater! Phew...", wz.GetName())
	return wz, nil
}
