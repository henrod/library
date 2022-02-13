package pg

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-pg/pg/v10"

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
		var pgErr pg.Error
		if errors.As(err, &pgErr) && pgErr.IntegrityViolation() {
			return nil, nil
		}

		return nil, fmt.Errorf("failed to insert shelf in postgres: %w", err)
	}

	return shelf.toEntity(), nil
}

func (g *Gateway) GetShelf(ctx context.Context, shelfName string) (*entities.Shelf, error) {
	shelf := new(Shelf)
	shelf.Name = shelfName

	err := g.db.ModelContext(ctx, shelf).
		WherePK().
		Select()
	if errors.Is(err, pg.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to select shelf in postgres: %w", err)
	}

	return shelf.toEntity(), nil
}
