package main

import (
	"context"
	"log"
	"net"

	pb "github.com/nazufel/telepresence-demo/proto/wizard"

	"google.golang.org/grpc"
)

const (
	port = ":9999"
)

type server struct {
	pb.UnimplementedWizardServiceServer
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("could not listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	pb.RegisterWizardServiceServer(grpcServer, &server{})

	log.Printf("running grpc server on port: %v", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to start the grpc server: %s", err)
	}
}

func (s *server) List(ctx context.Context, in *pb.Wizard) (*pb.Wizard, error) {
	if in.GetDeathEater() {
		log.Printf("Wizard: %s is a DeathEater! Run!", in.GetName())
		return in, nil
	}
	log.Printf("%s is not a DeathEater! Phew...", in.GetName())
	return in, nil
}
