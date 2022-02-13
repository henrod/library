package shelves

import (
	"context"
	"fmt"

	"github.com/Henrod/library/domain/errors"

	"github.com/Henrod/library/domain/entities"
)

type GetShelfDomain struct {
	gateway GetShelfGateway
}

type GetShelfGateway interface {
	GetShelf(ctx context.Context, shelfName string) (*entities.Shelf, error)
}

func NewGetShelfDomain(
	gateway GetShelfGateway,
) *GetShelfDomain {
	return &GetShelfDomain{
		gateway: gateway,
	}
}

func (g *GetShelfDomain) GetShelf(ctx context.Context, shelfName string) (*entities.Shelf, error) {
	shelf, err := g.gateway.GetShelf(ctx, shelfName)
	if err != nil {
		return nil, fmt.Errorf("failed to get shelf from gateway: %w", err)
	}

	if shelf == nil {
		return nil, errors.NotFoundError{
			Details: fmt.Sprintf("shelf %s not found", shelfName),
		}
	}

	return &entities.Shelf{
		Name:       shelf.Name,
		CreateTime: shelf.CreateTime,
		UpdateTime: shelf.UpdateTime,
	}, nil
}
