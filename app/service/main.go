package main

import (
	"fmt"
	"log"
	"net"

	proto "github.com/Henrod/library/protogen/go/api/v1"
	library "github.com/Henrod/library/service/api/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	listenOn := "127.0.0.1:8080"
	listener, err := net.Listen("tcp", listenOn)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", listenOn, err)
	}

	server := grpc.NewServer()
	proto.RegisterLibraryServiceServer(server, &library.LibraryService{})
	reflection.Register(server)
	log.Println("Listening on", listenOn)
	if err := server.Serve(listener); err != nil {
		return fmt.Errorf("failed to serve gRPC server: %w", err)
	}

	return nil
}
