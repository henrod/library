package pg

import (
	"context"
	"fmt"
	"time"

	"github.com/Henrod/library/domain/entities"
)

type Shelf struct {
	Name       string `pg:",pk"`
	CreateTime time.Time
	UpdateTime time.Time
}

func (s *Shelf) toEntity() *entities.Shelf {
	return &entities.Shelf{
		Name:       s.Name,
		CreateTime: s.CreateTime,
		UpdateTime: s.UpdateTime,
	}
}

func (g *Gateway) CreateShelf(ctx context.Context, eShelf *entities.Shelf) (*entities.Shelf, error) {
	now := time.Now()

	shelf := &Shelf{
		Name:       eShelf.Name,
		CreateTime: now,
		UpdateTime: now,
	}

	_, err := g.db.ModelContext(ctx, shelf).Insert()
	if err != nil {
		return nil, fmt.Errorf("failed to insert shelf in postgres: %w", err)
	}

	return shelf.toEntity(), nil
}
