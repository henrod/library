package books

import (
	"context"
	"fmt"

	"github.com/Henrod/library/domain/entities"
)

type ListBooksDomain struct {
	gateway ListBooksGateway
}

func NewListBooks(gateway ListBooksGateway) *ListBooksDomain {
	return &ListBooksDomain{gateway: gateway}
}

type ListBooksGateway interface {
	ListBooks(ctx context.Context, pageSize, pageOffset int) ([]*entities.Book, error)
	ListShelfBooks(ctx context.Context, shelfName string, pageSize, pageOffset int) ([]*entities.Book, error)
	CountBooks(ctx context.Context) (int, error)
	CountShelfBooks(ctx context.Context, shelfName string) (int, error)
}

func (l *ListBooksDomain) List(
	ctx context.Context,
	shelfName string,
	pageSize, pageOffset int,
) (books []*entities.Book, finished bool, err error) {
	var totalBooks int

	if shelfName == "-" {
		books, totalBooks, err = l.listBooks(ctx, pageSize, pageOffset)
		if err != nil {
			return nil, false, fmt.Errorf("failed to list books: %w", err)
		}
	} else {
		books, totalBooks, err = l.listShelfBooks(ctx, shelfName, pageSize, pageOffset)
		if err != nil {
			return nil, false, fmt.Errorf("failed to list shelf books: %w", err)
		}
	}

	finished = totalBooks <= pageOffset+pageSize

	return books, finished, nil
}

func (l *ListBooksDomain) listShelfBooks(
	ctx context.Context,
	shelfName string,
	pageSize, pageOffset int,
) (books []*entities.Book, totalBooks int, err error) {
	books, err = l.gateway.ListShelfBooks(ctx, shelfName, pageSize, pageOffset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list shelf books in gateway: %w", err)
	}

	totalBooks, err = l.gateway.CountShelfBooks(ctx, shelfName)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count shelf books in gateway: %w", err)
	}

	return books, totalBooks, nil
}

func (l *ListBooksDomain) listBooks(
	ctx context.Context,
	pageSize, pageOffset int,
) (books []*entities.Book, totalBooks int, err error) {
	books, err = l.gateway.ListBooks(ctx, pageSize, pageOffset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list shelf books in gateway: %w", err)
	}

	totalBooks, err = l.gateway.CountBooks(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count shelf books in gateway: %w", err)
	}

	return books, totalBooks, nil
}
