package api

import (
	"errors"

	domainErrors "github.com/Henrod/library/domain/errors"
	"github.com/golang/protobuf/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ToGRPCError(err error, details ...proto.Message) error {
	if err == nil {
		return nil
	}

	var grpcError error

	notFoundError := &domainErrors.NotFoundError{Details: ""}
	if errors.As(err, notFoundError) {
		grpcError = withDetails(codes.NotFound, "resource not found", details...)
	}

	alreadyExistsError := &domainErrors.AlreadyExistsError{Details: ""}
	if errors.As(err, alreadyExistsError) {
		grpcError = withDetails(codes.AlreadyExists, "resource already exists", details...)
	}

	if grpcError == nil {
		grpcError = status.Errorf(codes.Internal, "internal error")
	}

	return grpcError
}

func withDetails(code codes.Code, msg string, details ...proto.Message) error {
	errStatus := status.New(code, msg)
	detailedStatus, err := errStatus.WithDetails(details...)
	if err != nil {
		zap.S().With(zap.Error(err)).Error("failed to add details to error, skipping them")
		detailedStatus = errStatus
	}

	return detailedStatus.Err() //nolint:wrapcheck
}
