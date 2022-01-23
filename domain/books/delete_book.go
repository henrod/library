package books

import (
	"context"
	"fmt"

	"github.com/Henrod/library/domain/errors"
)

type DeleteBookDomain struct {
	gateway DeleteBookGateway
}

func NewDeleteBookDomain(gateway DeleteBookGateway) *DeleteBookDomain {
	return &DeleteBookDomain{gateway: gateway}
}

type DeleteBookGateway interface {
	DeleteBook(ctx context.Context, shelfName, bookName string) (bool, error)
}

func (g *DeleteBookDomain) DeleteBook(ctx context.Context, shelfName, bookName string) error {
	deleted, err := g.gateway.DeleteBook(ctx, shelfName, bookName)
	if err != nil {
		return fmt.Errorf("failed to get book from gateway: %w", err)
	}

	if !deleted {
		return errors.NotFoundError{
			Details: fmt.Sprintf("book %s at shelf %s not found", shelfName, bookName),
		}
	}

	return nil
}
