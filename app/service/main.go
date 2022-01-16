package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/Henrod/library/gateways/pg"

	"go.uber.org/zap"

	"github.com/Henrod/library/domain/books"

	proto "github.com/Henrod/library/protogen/go/api/v1"
	library "github.com/Henrod/library/service/api/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	ServerURL   = "127.0.0.1:8080"
	PostgresURL = "postgresql://postgres:password@localhost:5432/library?sslmode=disable"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	listener, err := net.Listen("tcp", ServerURL)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", ServerURL, err)
	}

	ctx := context.Background()

	logger, _ := zap.NewProduction()
	defer func() { _ = logger.Sync() }()
	sugar := logger.Sugar()

	gateway, err := pg.NewGateway(ctx, PostgresURL)
	if err != nil {
		return fmt.Errorf("failed to create gateway: %w", err)
	}

	server := grpc.NewServer()
	reflection.Register(server)

	proto.RegisterLibraryServiceServer(server, library.NewLibraryService(
		books.NewListBooks(gateway),
		sugar,
	))

	sugar.Infof("Listening on %s", ServerURL)
	if err := server.Serve(listener); err != nil {
		return fmt.Errorf("failed to serve gRPC server: %w", err)
	}

	return nil
}
