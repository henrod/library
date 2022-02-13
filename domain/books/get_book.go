package books

import (
	"context"
	"fmt"

	"github.com/Henrod/library/domain/errors"

	"github.com/Henrod/library/domain/entities"
)

type GetBookDomain struct {
	gateway GetBookGateway
}

func NewGetBookDomain(gateway GetBookGateway) *GetBookDomain {
	return &GetBookDomain{gateway: gateway}
}

type GetBookGateway interface {
	GetBook(ctx context.Context, shelfName, bookName string) (*entities.Book, error)
}

func (g *GetBookDomain) GetBook(ctx context.Context, shelfName, bookName string) (*entities.Book, error) {
	book, err := g.gateway.GetBook(ctx, shelfName, bookName)
	if err != nil {
		return nil, fmt.Errorf("failed to get book from gateway: %w", err)
	}

	if book == nil {
		return nil, errors.NotFoundError{
			Details: fmt.Sprintf("book %s at shelf %s not found", shelfName, bookName),
		}
	}

	return &entities.Book{
		Name:       book.Name,
		Author:     book.Author,
		Shelf:      book.Shelf,
		CreateTime: book.CreateTime,
		UpdateTime: book.UpdateTime,
	}, nil
}
