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

func (c *CreateBookDomain) CreateBook(
	ctx context.Context,
	shelfName string,
	inputBook *entities.Book,
) (*entities.Book, error) {
	bookName := inputBook.Name
	book, err := c.gateway.CreateBook(ctx, shelfName, inputBook)
	if err != nil {
		return nil, fmt.Errorf("failed to create book in gateway: %w", err)
	}

	if book == nil {
		return nil, errors.AlreadyExistsError{
			Details: fmt.Sprintf("book %s at shelf %s already exists", bookName, shelfName),
		}
	}

	return book, nil
}
