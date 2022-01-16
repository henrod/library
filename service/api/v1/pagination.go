package v1

import (
	"encoding/base64"
	"strconv"

	"go.uber.org/zap"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const defaultPageSize = 10

func errInvalidPageToken() error {
	errStatus, _ := status.New(codes.InvalidArgument, "invalid page token").WithDetails(&errdetails.BadRequest{
		FieldViolations: []*errdetails.BadRequest_FieldViolation{{
			Field:       "page_token",
			Description: "page_token must be a previously returned next_page_token",
		}},
	})

	return errStatus.Err() //nolint:wrapcheck
}

func getPageOffset(log *zap.SugaredLogger, pageToken string) (int, error) {
	if pageToken == "" {
		return 0, nil
	}

	bPageToken, err := base64.StdEncoding.DecodeString(pageToken)
	if err != nil {
		log.With(
			zap.Error(err),
			zap.String("page_token", pageToken),
		).Error("failed to decode base64 page_token")

		return 0, errInvalidPageToken()
	}

	pageOffset, err := strconv.Atoi(string(bPageToken))
	if err != nil {
		log.With(
			zap.Error(err),
			zap.ByteString("page_token", bPageToken),
		).Error("failed to decode base64 page_token")

		return 0, errInvalidPageToken()
	}

	return pageOffset, nil
}

func getPageSize(pageSize int32) int {
	if pageSize == 0 {
		return defaultPageSize
	}

	return int(pageSize)
}

func getNextIntPageToken(next int) string {
	bNext := []byte(strconv.Itoa(next))

	return base64.StdEncoding.EncodeToString(bNext)
}
