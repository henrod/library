package pg

import (
	"context"
	"fmt"
	"time"

	"github.com/Henrod/library/domain/entities"
)

type Shelf struct {
	Name string `pg:",pk"`
}

type Book struct {
	ShelfName  string
	Shelf      *Shelf `pg:"rel:has-one"`
	Author     string `pg:",pk"`
	Title      string `pg:",pk"`
	CreateTime time.Time
	UpdateTime time.Time
}

func (g *Gateway) ListBooks(ctx context.Context, shelf string) ([]*entities.Book, error) {
	var books []*Book
	err := g.db.ModelContext(ctx, &books).
		Relation("Shelf").
		Where("book.shelf_name = ?", shelf).
		Select()
	if err != nil {
		return nil, fmt.Errorf("failed to select books from shelf in postgres: %w", err)
	}

	eBooks := make([]*entities.Book, len(books))
	for i, book := range books {
		eBooks[i] = &entities.Book{
			Title:      book.Title,
			Author:     book.Author,
			CreateTime: book.CreateTime,
			UpdateTime: book.UpdateTime,
		}
	}

	return eBooks, nil
}
