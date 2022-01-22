package api

import (
	"errors"

	domainErrors "github.com/Henrod/library/domain/errors"
	"google.golang.org/genproto/googleapis/rpc/errdetails"

	// TODO: fix this linter error: github.com/golang/protobuf/proto incompatible with google.golang.org/protobuf/proto
	"github.com/golang/protobuf/proto" //nolint:staticcheck
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Details map[codes.Code][]proto.Message

func GRPCError(err error, details Details) error {
	if err == nil {
		return nil
	}

	var grpcError error

	notFoundError := new(domainErrors.NotFoundError)
	if errors.As(err, notFoundError) {
		grpcError = withDetails(codes.NotFound, "resource not found", details[codes.NotFound]...)
	}

	alreadyExistsError := new(domainErrors.AlreadyExistsError)
	if errors.As(err, alreadyExistsError) {
		grpcError = withDetails(codes.AlreadyExists, "resource already exists", details[codes.AlreadyExists]...)
	}

	badRequestError := new(domainErrors.BadRequestError)
	if errors.As(err, badRequestError) {
		grpcError = withDetails(codes.InvalidArgument, "invalid argument", details[codes.InvalidArgument]...)
	}

	if grpcError == nil {
		grpcError = status.Errorf(codes.Internal, "internal error")
	}

	return grpcError
}

func BadRequestDetails(err error) (*errdetails.BadRequest, bool) {
	if err == nil {
		return nil, false
	}

	var badRequestError *domainErrors.BadRequestError
	if !errors.As(err, &badRequestError) {
		return nil, false
	}

	return &errdetails.BadRequest{
		FieldViolations: []*errdetails.BadRequest_FieldViolation{
			{
				Field:       badRequestError.InvalidField,
				Description: badRequestError.Details,
			},
		},
	}, true
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
