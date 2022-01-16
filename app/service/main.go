package main

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/Henrod/library/domain/books"
	"github.com/Henrod/library/gateways/pg"
	proto "github.com/Henrod/library/protogen/go/api/v1"
	library "github.com/Henrod/library/service/api/v1"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

const (
	GRPCServerURL = "127.0.0.1:8080"
	HTTPServerURL = "127.0.0.1:8081"
	PostgresURL   = "postgresql://postgres:password@localhost:5432/library?sslmode=disable"
)

func main() {
	logger, _ := zap.NewProduction()
	defer func() { _ = logger.Sync() }()
	sugar := logger.Sugar()

	if err := run(sugar); err != nil {
		sugar.Fatal(err)
	}
}

func runHTTPServer(ctx context.Context, sugar *zap.SugaredLogger) error {
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err := proto.RegisterLibraryServiceHandlerFromEndpoint(ctx, mux, GRPCServerURL, opts)
	if err != nil {
		return fmt.Errorf("failed to connect to gRPC service: %w", err)
	}

	sugar.Infof("Listening HTTP on %s", HTTPServerURL)
	err = http.ListenAndServe(HTTPServerURL, mux)
	if err != nil {
		return fmt.Errorf("failed to listen and serve HTTP server: %w", err)
	}

	return nil
}

func run(sugar *zap.SugaredLogger) error {
	listener, err := net.Listen("tcp", GRPCServerURL)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", GRPCServerURL, err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	gateway, err := pg.NewGateway(ctx, PostgresURL)
	if err != nil {
		return fmt.Errorf("failed to create gateway: %w", err)
	}

	server := grpc.NewServer()
	reflection.Register(server)

	proto.RegisterLibraryServiceServer(server, library.NewLibraryService(
		sugar,
		books.NewListBooks(gateway),
		books.NewGetBookDomain(gateway),
	))

	go func() {
		if err = runHTTPServer(ctx, sugar); err != nil {
			sugar.Fatal(err)
		}
	}()

	sugar.Infof("Listening gRPC on %s", GRPCServerURL)
	if err := server.Serve(listener); err != nil {
		return fmt.Errorf("failed to serve gRPC server: %w", err)
	}

	return nil
}
