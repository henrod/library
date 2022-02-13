package v1

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Henrod/library/domain/books"
	"github.com/Henrod/library/domain/entities"
	"github.com/Henrod/library/domain/shelves"
	v1 "github.com/Henrod/library/protogen/go/api/v1"
	"github.com/Henrod/library/service/api"
	"go.uber.org/zap"
	"google.golang.org/genproto/googleapis/longrunning"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	// TODO: fix this linter error: github.com/golang/protobuf/proto incompatible with google.golang.org/protobuf/proto.
	"github.com/golang/protobuf/proto" //nolint:staticcheck
)

type LibraryService struct {
	listBooks   *books.ListBooksDomain
	getBook     *books.GetBookDomain
	createBook  *books.CreateBookDomain
	updateBook  *books.UpdateBookDomain
	deleteBook  *books.DeleteBookDomain
	getShelf    *shelves.GetShelfDomain
	createShelf *shelves.CreateShelfDomain
	log         *zap.SugaredLogger
}

func NewLibraryService(
	log *zap.SugaredLogger,
	listBooks *books.ListBooksDomain,
	getBook *books.GetBookDomain,
	createBook *books.CreateBookDomain,
	updateBook *books.UpdateBookDomain,
	deleteBook *books.DeleteBookDomain,
	getShelf *shelves.GetShelfDomain,
	createShelf *shelves.CreateShelfDomain,
) *LibraryService {
	return &LibraryService{
		log:         log,
		listBooks:   listBooks,
		getBook:     getBook,
		createBook:  createBook,
		updateBook:  updateBook,
		deleteBook:  deleteBook,
		getShelf:    getShelf,
		createShelf: createShelf,
	}
}

// ListBooks returns a list of books in the following cases:
// 1. When shelfID is valid, returns the books in that shelf.
// 2. When shelfID is "-", returns all books in the library.
// 3. When shelfID is invalid, returns empty list.
//
// Method is paginated in the following standard: https://cloud.google.com/apis/design/design_patterns#list_pagination.
//
// Book resource name must be in the format: /shelves/:shelfID/books.
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

		return nil, api.GRPCError(err, nil) //nolint:wrapcheck
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

		return nil, api.GRPCError(err, api.Details{ //nolint:wrapcheck
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
	inputBook := &entities.Book{
		Name:       request.GetBook().GetName(),
		Author:     request.GetBook().GetAuthor(),
		Shelf:      nil,
		CreateTime: time.Time{},
		UpdateTime: time.Time{},
	}

	book, err := l.createBook.CreateBook(ctx, shelfName, inputBook)
	if err != nil {
		l.log.With(zap.Error(err)).Error("failed to create book in domain")

		return nil, api.GRPCError(err, api.Details{ //nolint:wrapcheck
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

func (l *LibraryService) UpdateBook(ctx context.Context, request *v1.UpdateBookRequest) (*v1.Book, error) {
	name := strings.Split(request.GetBook().GetName(), "/")
	if len(name) != 4 {
		err := status.Errorf(codes.InvalidArgument, "book name must be of format 'shelves/*/books/*'")

		return nil, fmt.Errorf("failed to get book name: %w", err)
	}

	shelfName := name[1]
	inputBook := &entities.Book{
		Name:       name[3],
		Author:     request.GetBook().GetAuthor(),
		Shelf:      nil,
		CreateTime: time.Time{},
		UpdateTime: time.Time{},
	}

	book, err := l.updateBook.UpdateBook(ctx, shelfName, inputBook, request.GetUpdateMask())
	if err != nil {
		l.log.With(zap.Error(err)).Error("failed to update book in domain")

		details := api.Details{
			codes.NotFound: {&errdetails.ResourceInfo{
				ResourceType: "book",
				ResourceName: request.GetBook().GetName(),
				Owner:        fmt.Sprintf("%s/%s", name[0], name[1]),
				Description:  "book not found in shelf",
			}},
		}

		if badRequestDetail, ok := api.BadRequestDetails(err); ok {
			details[codes.InvalidArgument] = []proto.Message{badRequestDetail}
		}

		return nil, api.GRPCError(err, details) //nolint:wrapcheck
	}

	return toProtoBook(book), nil
}

func (l *LibraryService) DeleteBook(ctx context.Context, request *v1.DeleteBookRequest) (*emptypb.Empty, error) {
	name := strings.Split(request.GetName(), "/")
	if len(name) != 4 {
		err := status.Errorf(codes.InvalidArgument, "book name must be of format 'shelves/*/books/*'")

		return nil, fmt.Errorf("failed to get book name: %w", err)
	}

	shelfName := name[1]
	bookName := name[3]

	err := l.deleteBook.DeleteBook(ctx, shelfName, bookName)
	if err != nil {
		l.log.With(zap.Error(err)).Error("failed to update book in domain")

		return nil, api.GRPCError(err, api.Details{ //nolint:wrapcheck
			codes.NotFound: {&errdetails.ResourceInfo{
				ResourceType: "book",
				ResourceName: request.GetName(),
				Owner:        fmt.Sprintf("%s/%s", name[0], name[1]),
				Description:  "book not found in shelf",
			}},
		})
	}

	return &emptypb.Empty{}, nil
}

func (l *LibraryService) CreateShelf(
	ctx context.Context,
	request *v1.CreateShelfRequest,
) (*longrunning.Operation, error) {
	inputShelf := &entities.Shelf{
		Name:       request.GetShelf().GetName(),
		CreateTime: time.Time{},
		UpdateTime: time.Time{},
	}

	operation, err := l.createShelf.StartCreateShelfOperation(ctx, inputShelf)
	if err != nil {
		l.log.With(zap.Error(err)).Error("failed to start create shelf operation in domain")

		return nil, api.GRPCError(err, api.Details{ //nolint:wrapcheck
			codes.AlreadyExists: {&errdetails.ResourceInfo{
				ResourceType: "operation",
				ResourceName: l.createShelf.GetOperationName(inputShelf.Name),
				Owner:        "library",
				Description:  "the create shelf operation already exists",
			}},
		})
	}

	return toLongRunningOperation("CreateShelf", operation), nil
}

func (l *LibraryService) GetOperation(
	ctx context.Context,
	request *v1.GetOperationRequest,
) (*longrunning.Operation, error) {
	parts := strings.Split(request.GetName(), "/")
	// assume only long-running operation is CreateShelf.
	// request.name = operations/shelves/{shelf_name}
	if len(parts) != 3 {
		err := status.Errorf(codes.InvalidArgument, "operation name must be of format 'operations/shelves/*'")

		return nil, fmt.Errorf("failed to get operation: %w", err)
	}

	shelfName := parts[len(parts)-1]
	longRunningOperationName := "CreateShelf"

	operation, err := l.createShelf.GetOperation(shelfName)
	if err != nil {
		return nil, api.GRPCError(err, api.Details{ //nolint:wrapcheck
			codes.NotFound: {&errdetails.ResourceInfo{
				ResourceType: "operation",
				ResourceName: l.createShelf.GetOperationName(shelfName),
				Owner:        "library",
				Description:  "the operation doesn't exist; is not running nor completed",
			}},
		})
	}

	longRunningOperation := toLongRunningOperation(longRunningOperationName, operation)
	if !longRunningOperation.Done || longRunningOperation.Result != nil {
		return longRunningOperation, nil
	}

	shelf, err := l.getShelf.GetShelf(ctx, shelfName)
	if err != nil {
		l.log.With(zap.Error(err)).Error("failed to get shelf in domain")

		return nil, api.GRPCError(err, api.Details{ //nolint:wrapcheck
			codes.NotFound: {&errdetails.ResourceInfo{
				ResourceType: "operation",
				ResourceName: request.GetName(),
				Owner:        "library",
				Description:  "the operation doesn't exist; is not running nor completed",
			}},
		})
	}

	response, _ := anypb.New(toProtoShelf(shelf))
	longRunningOperation.Result = &longrunning.Operation_Response{Response: response}
	longRunningOperation.Done = true

	return longRunningOperation, nil
}

func toProtoBook(book *entities.Book) *v1.Book {
	return &v1.Book{
		Name:       bookResourceName(book),
		Author:     book.Author,
		CreateTime: timestamppb.New(book.CreateTime),
		UpdateTime: timestamppb.New(book.UpdateTime),
	}
}

func bookResourceName(book *entities.Book) string {
	return fmt.Sprintf("%s/books/%s", shelfResourceName(book.Shelf), book.Name)
}

func toProtoShelf(shelf *entities.Shelf) *v1.Shelf {
	if shelf == nil {
		return nil
	}

	return &v1.Shelf{
		Name:       shelfResourceName(shelf),
		CreateTime: timestamppb.New(shelf.CreateTime),
		UpdateTime: timestamppb.New(shelf.UpdateTime),
	}
}

func shelfResourceName(shelf *entities.Shelf) string {
	return fmt.Sprintf("shelves/%s", shelf.Name)
}

func toLongRunningOperation(name string, operation *entities.Operation) *longrunning.Operation {
	metadata, _ := anypb.New(&v1.Operation{
		Name:       name,
		Stage:      operation.Stage,
		Percentage: uint32(operation.Percentage),
	})

	longRunningOperation := &longrunning.Operation{
		Name:     operation.Name,
		Metadata: metadata,
		Done:     operation.Finished(),
		Result:   nil,
	}

	if operation.Error != nil {
		err := api.GRPCError(operation.Error, api.Details{
			codes.AlreadyExists: {&errdetails.ResourceInfo{
				ResourceType: "shelf",
				ResourceName: operation.ResourceName(),
				Owner:        "library",
				Description:  "the shelf already exists in the library",
			}},
		})

		longRunningOperation.Result = &longrunning.Operation_Error{
			Error: status.Convert(err).Proto(),
		}
	}

	return longRunningOperation
}
