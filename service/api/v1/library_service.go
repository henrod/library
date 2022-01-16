package v1

import (
	"context"

	v1 "github.com/Henrod/library/protogen/go/api/v1"
)

type LibraryService struct{}

func (l *LibraryService) ListBooks(ctx context.Context, request *v1.ListBooksRequest) (*v1.ListBooksResponse, error) {
	return &v1.ListBooksResponse{
		Books:         nil,
		NextPageToken: "",
	}, nil
}
