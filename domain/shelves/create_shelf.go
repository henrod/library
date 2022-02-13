package shelves

import (
	"context"
	"fmt"

	"github.com/Henrod/library/domain/entities"
	"github.com/Henrod/library/domain/errors"
)

type CreateShelfDomain struct {
	gateway CreateShelfGateway
}

func NewCreateShelfDomain(gateway CreateShelfGateway) *CreateShelfDomain {
	return &CreateShelfDomain{gateway: gateway}
}

type CreateShelfGateway interface {
	CreateShelf(ctx context.Context, shelf *entities.Shelf) (*entities.Shelf, error)
}

func (c *CreateShelfDomain) CreateShelf(
	ctx context.Context,
	shelf *entities.Shelf,
) (*entities.Shelf, error) {
	book, err := c.gateway.CreateShelf(ctx, shelf)
	if err != nil {
		return nil, fmt.Errorf("failed to create book in gateway: %w", err)
	}

	if book == nil {
		return nil, errors.AlreadyExistsError{
			Details: fmt.Sprintf("shelf %s already exists", shelf.Name),
		}
	}

	return book, nil
}
