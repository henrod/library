package v1

import (
	"context"
	"fmt"
	"strings"

	"github.com/Henrod/library/service/api"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/Henrod/library/domain/books"
	"go.uber.org/zap"

	v1 "github.com/Henrod/library/protogen/go/api/v1"
)

type LibraryService struct {
	listBooks *books.ListBooksDomain
	getBook   *books.GetBookDomain
	log       *zap.SugaredLogger
}

func NewLibraryService(
	log *zap.SugaredLogger,
	listBooks *books.ListBooksDomain,
	getBook *books.GetBookDomain,
) *LibraryService {
	return &LibraryService{
		log:       log,
		listBooks: listBooks,
		getBook:   getBook,
	}
}

func (l *LibraryService) ListBooks(ctx context.Context, request *v1.ListBooksRequest) (*v1.ListBooksResponse, error) {
	parent := strings.Split(request.GetParent(), "/")
	if len(parent) != 2 {
		err := status.Errorf(codes.InvalidArgument, "parent must be of format 'shelves/*'")

		return nil, fmt.Errorf("failed to get parent: %w", err)
	}

	shelfName := parent[1]
	pageSize := getPageSize(request.GetPageSize())
	pageOffset, err := getPageOffset(l.log, request.GetPageToken())
	if err != nil {
		return nil, err
	}

	eBooks, finishedBooks, err := l.listBooks.List(ctx, shelfName, pageSize, pageOffset)
	if err != nil {
		l.log.With(zap.Error(err)).Error("failed to list books in domain")

		return nil, api.ToGRPCError(err) //nolint:wrapcheck
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

func (l *LibraryService) GetBook(ctx context.Context, request *v1.GetBookRequest) (*v1.Book, error) {
	resourceName := strings.Split(request.GetName(), "/")
	if len(resourceName) != 4 {
		err := status.Errorf(codes.InvalidArgument, "resource name must be of format 'shelves/*/books/*'")

		return nil, fmt.Errorf("failed to get resource name: %w", err)
	}

	shelfName := resourceName[1]
	bookName := resourceName[3]

	book, err := l.getBook.GetBook(ctx, shelfName, bookName)
	if err != nil {
		l.log.With(zap.Error(err)).Error("failed to get book in domain")

		return nil, api.ToGRPCError(err) //nolint:wrapcheck
	}

	return &v1.Book{
		Isbn:       book.ISBN,
		Title:      book.Title,
		Author:     book.Author,
		CreateTime: timestamppb.New(book.CreateTime),
		UpdateTime: timestamppb.New(book.UpdateTime),
	}, nil
}
