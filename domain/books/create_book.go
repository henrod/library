package books

import (
	"context"
	"fmt"

	"github.com/Henrod/library/domain/entities"
	"github.com/Henrod/library/domain/errors"
)

type CreateBookDomain struct {
	gateway CreateBookGateway
}

func NewCreateBookDomain(gateway CreateBookGateway) *CreateBookDomain {
	return &CreateBookDomain{gateway: gateway}
}

type CreateBookGateway interface {
	CreateBook(ctx context.Context, shelfName string, book *entities.Book) (*entities.Book, error)
}

func (g *CreateBookDomain) CreateBook(
	ctx context.Context,
	shelfName string,
	book *entities.Book,
) (*entities.Book, error) {
	bookName := book.Name
	book, err := g.gateway.CreateBook(ctx, shelfName, book)
	if err != nil {
		return nil, fmt.Errorf("failed to create book in gateway: %w", err)
	}

	if book == nil {
		err := errors.AlreadyExistsError{
			Details: fmt.Sprintf("book %s at shelf %s already exists", bookName, shelfName),
		}

		return nil, err
	}

	return &entities.Book{
		Name:       book.Name,
		Author:     book.Author,
		CreateTime: book.CreateTime,
		UpdateTime: book.UpdateTime,
	}, nil
}
