package main

import (
	"context"
	"log"
	"net"
	"os"
	"time"

	pb "github.com/nazufel/telepresence-demo/proto/wizard"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
)

const (
	grpcPort = ":9999"
)

var (
	dbName       string
	dbConnection string
	dbHost       string
	dbUsername   string
	dbPassword   string
	ctx          context.Context
)

type server struct {
	pb.UnimplementedWizardServiceServer
}

func main() {
	log.Println("starting wizards server...")
	log.Println("trying to connect to database")

	_, err := newStorage()
	if err != nil {
		log.Fatalf("error setting up new storage: %v", err)
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

// Storage struct holds information about connecting to the DB
type storage struct {
	db *mongo.Database
	c  *mongo.Client
}

func newStorage() (*storage, error) {

	var err error

	if len(os.Getenv("MONGO_DATABASE")) == 0 {
		log.Fatal("MONGO_DATABASE is a required env variable. Please define a database to use.")
	}
	dbName = os.Getenv("MONGO_DATABASE")

	if len(os.Getenv("MONGO_HOST")) == 0 {
		log.Fatal("MONGO_HOST is a required env variable. Please define a database host string to use.")
	}
	dbHost = os.Getenv("MONGO_HOST")

	if len(os.Getenv("MONGO_USER")) == 0 {
		log.Fatal("MONGO_USER is a required env variable. Please define a database username string to use.")
	}
	dbUsername = os.Getenv("MONGO_USER")

	if len(os.Getenv("MONGO_PASSWORD")) == 0 {
		log.Fatal("MONGO_PASSWORD is a required env variable. Please define a database password string to use.")
	}
	dbPassword = os.Getenv("MONGO_PASSWORD")

	s := new(storage)

	dbConnection = "mongodb://" + dbUsername + ":" + dbPassword + "@" + dbHost + ":27017" + "/" + dbName + "?authSource=admin"

	client, err := mongo.NewClient(options.Client().ApplyURI(dbConnection))
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			log.Fatalf("unable to disconnect client: %v", err)
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Printf("attempting to establish database connection: %s", dbConnection)

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal("unable to connect to database: %s %v", dbConnection, err)
	}
	s.db = client.Database(dbName)
	s.c = client

	log.Printf("database connection successfully established: %s", dbConnection)

	return s, nil
}

func (s *server) List(ctx context.Context, in *pb.Wizard) (*pb.Wizard, error) {
	if in.GetDeathEater() {
		log.Printf("Wizard: %s is a DeathEater! Run!", in.GetName())
		return in, nil
	}
	log.Printf("%s is not a DeathEater! Phew...", in.GetName())
	return in, nil
}
