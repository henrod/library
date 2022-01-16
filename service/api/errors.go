package api

import (
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	domainErrors "github.com/Henrod/library/domain/errors"
)

func ToGRPCError(err error) error {
	if err == nil {
		return nil
	}

	var grpcError error

	notFoundError := &domainErrors.NotFoundError{Details: ""}
	if errors.As(err, notFoundError) {
		grpcError = status.Error(codes.NotFound, "resource not found")
	}

	if grpcError == nil {
		grpcError = status.Errorf(codes.Internal, "internal error")
	}

	return grpcError
}
