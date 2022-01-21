package v1

import (
	"context"
	"fmt"
	"strings"
	"time"

	"google.golang.org/genproto/googleapis/rpc/errdetails"

	"github.com/Henrod/library/domain/entities"

	"github.com/Henrod/library/service/api"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/Henrod/library/domain/books"
	"go.uber.org/zap"

	v1 "github.com/Henrod/library/protogen/go/api/v1"
)

type LibraryService struct {
	listBooks  *books.ListBooksDomain
	getBook    *books.GetBookDomain
	createBook *books.CreateBookDomain
	log        *zap.SugaredLogger
}

func NewLibraryService(
	log *zap.SugaredLogger,
	listBooks *books.ListBooksDomain,
	getBook *books.GetBookDomain,
	createBook *books.CreateBookDomain,
) *LibraryService {
	return &LibraryService{
		log:        log,
		listBooks:  listBooks,
		getBook:    getBook,
		createBook: createBook,
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

		return nil, api.ToGRPCError(err, nil) //nolint:wrapcheck
	}

	nextPageToken := ""
	if !finishedBooks {
		nextPageToken = getNextIntPageToken(pageOffset + 1)
	}

	pBooks := make([]*v1.Book, len(eBooks))
	for i, book := range eBooks {
		pBooks[i] = toProtoBook(book)
	}

	return &v1.ListBooksResponse{
		Books:         pBooks,
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

		return nil, api.ToGRPCError(err, api.Details{ //nolint:wrapcheck
			codes.NotFound: {&errdetails.ResourceInfo{
				ResourceType: "book",
				ResourceName: request.GetName(),
				Owner:        shelfName,
				Description:  "the book does not exist in shelf",
			}},
		})
	}

	return toProtoBook(book), nil
}

func (l *LibraryService) CreateBook(ctx context.Context, request *v1.CreateBookRequest) (*v1.Book, error) {
	parent := strings.Split(request.GetParent(), "/")
	if len(parent) != 2 {
		err := status.Errorf(codes.InvalidArgument, "parent must be of format 'shelves/*'")

		return nil, fmt.Errorf("failed to get parent: %w", err)
	}

	shelfName := parent[1]
	book := &entities.Book{
		Name:       request.GetBook().GetName(),
		Author:     request.GetBook().GetAuthor(),
		CreateTime: time.Time{},
		UpdateTime: time.Time{},
	}

	book, err := l.createBook.CreateBook(ctx, shelfName, book)
	if err != nil {
		l.log.With(zap.Error(err)).Error("failed to create book in domain")

		return nil, api.ToGRPCError(err, api.Details{ //nolint:wrapcheck
			codes.AlreadyExists: {&errdetails.ResourceInfo{
				ResourceType: "book",
				ResourceName: request.GetBook().GetName(),
				Owner:        request.GetParent(),
				Description:  "the book already exists in shelf",
			}},
		})
	}

	return toProtoBook(book), nil
}

func toProtoBook(book *entities.Book) *v1.Book {
	return &v1.Book{
		Name:       book.Name,
		Author:     book.Author,
		CreateTime: timestamppb.New(book.CreateTime),
		UpdateTime: timestamppb.New(book.UpdateTime),
	}
}
