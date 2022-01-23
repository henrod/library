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
	Name string `pg:",pk"`
}

type Book struct {
	ShelfName  string `pg:",pk"`
	Shelf      *Shelf `pg:"rel:has-one"`
	Name       string `pg:",pk"`
	Author     string
	CreateTime time.Time
	UpdateTime time.Time
}

func (b *Book) toEntitiesBook() *entities.Book {
	return &entities.Book{
		Name:       b.Name,
		Author:     b.Author,
		CreateTime: b.CreateTime,
		UpdateTime: b.UpdateTime,
	}
}

func (g *Gateway) ListBooks(ctx context.Context, shelfName string, pageSize, pageOffset int) ([]*entities.Book, error) {
	var books []*Book
	err := g.db.ModelContext(ctx, &books).
		Relation("Shelf").
		Where("book.shelf_name = ?", shelfName).
		Limit(pageSize).
		Offset(pageOffset).
		Select()
	if err != nil {
		return nil, fmt.Errorf("failed to select books from shelf in postgres: %w", err)
	}

	eBooks := make([]*entities.Book, len(books))
	for i, book := range books {
		eBooks[i] = book.toEntitiesBook()
	}

	return eBooks, nil
}

func (g *Gateway) CountBooks(ctx context.Context, shelfName string) (int, error) {
	book := new(Book)
	count, err := g.db.ModelContext(ctx, book).
		Relation("Shelf").
		Where("book.shelf_name = ?", shelfName).
		Count()
	if err != nil {
		return 0, fmt.Errorf("failed to count books from shelf in postgres: %w", err)
	}

	return count, nil
}

// GetBook returns book of name from shelf.
// If book not found, returns nil book and nil error.
func (g *Gateway) GetBook(ctx context.Context, shelfName, bookName string) (*entities.Book, error) {
	book := new(Book)
	err := g.db.ModelContext(ctx, book).
		Relation("Shelf").
		Where("book.shelf_name = ?", shelfName).
		Where("book.name = ?", bookName).
		Select()
	if errors.Is(err, pg.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to select book in postgres: %w", err)
	}

	return book.toEntitiesBook(), nil
}

func (g *Gateway) CreateBook(ctx context.Context, shelfName string, eBook *entities.Book) (*entities.Book, error) {
	now := time.Now()

	book := &Book{
		ShelfName:  shelfName,
		Shelf:      nil,
		Name:       eBook.Name,
		Author:     eBook.Author,
		CreateTime: now,
		UpdateTime: now,
	}

	_, err := g.db.ModelContext(ctx, book).Insert()
	if err != nil {
		var pgErr pg.Error
		// TODO: fix when shelf doesn't exist
		if errors.As(err, &pgErr) && pgErr.IntegrityViolation() {
			return nil, nil
		}

		return nil, fmt.Errorf("failed to insert book in postgres: %w", err)
	}

	return book.toEntitiesBook(), nil
}

func (g *Gateway) UpdateBook(
	ctx context.Context,
	shelfName string,
	eBook *entities.Book,
	fields []string,
) (*entities.Book, error) {
	now := time.Now()

	book := &Book{
		ShelfName:  shelfName,
		Shelf:      nil,
		Name:       eBook.Name,
		Author:     eBook.Author,
		CreateTime: time.Time{},
		UpdateTime: now,
	}

	fields = append(fields, "update_time")

	_, err := g.db.ModelContext(ctx, book).Column(fields...).WherePK().Returning("*").Update()
	if err != nil {
		if errors.Is(err, pg.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("failed to insert book in postgres: %w", err)
	}

	// TODO: handle when shelf doesn't exist

	return book.toEntitiesBook(), nil
}
func (g *Gateway) DeleteBook(ctx context.Context, shelfName, bookName string) (bool, error) {
	book := &Book{ShelfName: shelfName, Name: bookName} //nolint:exhaustivestruct
	r, err := g.db.ModelContext(ctx, book).WherePK().Delete()
	if err != nil {
		if errors.Is(err, pg.ErrNoRows) {
			return false, nil
		}

		return false, fmt.Errorf("failed to delete book in postgres: %w", err)
	}

	deleted := r.RowsAffected() > 0

	return deleted, nil
}
