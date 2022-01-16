package v1

import (
	"context"
	"fmt"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/Henrod/library/domain/books"
	"go.uber.org/zap"

	v1 "github.com/Henrod/library/protogen/go/api/v1"
)

type LibraryService struct {
	listBooks *books.ListBooksDomain
	log       *zap.SugaredLogger
}

func NewLibraryService(listBooks *books.ListBooksDomain, log *zap.SugaredLogger) *LibraryService {
	return &LibraryService{listBooks: listBooks, log: log}
}

func (l *LibraryService) ListBooks(ctx context.Context, request *v1.ListBooksRequest) (*v1.ListBooksResponse, error) {
	parent := strings.Split(request.GetParent(), "/")
	if len(parent) != 2 {
		err := status.Errorf(codes.InvalidArgument, "parent must be of format 'parents/*'")

		return nil, fmt.Errorf("failed to split parent: %w", err)
	}

	shelf := parent[1]
	pageSize := getPageSize(request.GetPageSize())
	pageOffset, err := getPageOffset(l.log, request.GetPageToken())
	if err != nil {
		return nil, err
	}

	eBooks, finishedBooks, err := l.listBooks.List(ctx, shelf, pageSize, pageOffset)
	if err != nil {
		l.log.With(zap.Error(err)).Error("failed to list books in domain")

		return nil, fmt.Errorf("api error: %w", err)
	}

	nextPageToken := ""
	if !finishedBooks {
		nextPageToken = getNextIntPageToken(pageOffset + 1)
	}

	rBooks := make([]*v1.Book, len(eBooks))
	for i, eBook := range eBooks {
		rBooks[i] = &v1.Book{
			Isbn:       eBook.ISBN,
			Title:      eBook.Title,
			Author:     eBook.Author,
			CreateTime: timestamppb.New(eBook.CreateTime),
			UpdateTime: timestamppb.New(eBook.UpdateTime),
		}
	}

	return &v1.ListBooksResponse{
		Books:         rBooks,
		NextPageToken: nextPageToken,
	}, nil
}
