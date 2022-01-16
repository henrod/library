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
	ISBN       string `pg:",pk"`
	Author     string
	Title      string
	CreateTime time.Time
	UpdateTime time.Time
}

func (g *Gateway) ListBooks(ctx context.Context, shelf string, pageSize, pageOffset int) ([]*entities.Book, error) {
	var books []*Book
	err := g.db.ModelContext(ctx, &books).
		Relation("Shelf").
		Where("book.shelf_name = ?", shelf).
		Limit(pageSize).
		Offset(pageOffset).
		Select()
	if err != nil {
		return nil, fmt.Errorf("failed to select books from shelf in postgres: %w", err)
	}

	eBooks := make([]*entities.Book, len(books))
	for i, book := range books {
		eBooks[i] = &entities.Book{
			ISBN:       book.ISBN,
			Title:      book.Title,
			Author:     book.Author,
			CreateTime: book.CreateTime,
			UpdateTime: book.UpdateTime,
		}
	}

	return eBooks, nil
}

func (g *Gateway) CountBooks(ctx context.Context, shelf string) (int, error) {
	book := new(Book)
	count, err := g.db.ModelContext(ctx, book).
		Relation("Shelf").
		Where("book.shelf_name = ?", shelf).
		Count()
	if err != nil {
		return 0, fmt.Errorf("failed to count books from shelf in postgres: %w", err)
	}

	return count, nil
}
