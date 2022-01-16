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
	ListBooks(ctx context.Context, shelf string, pageSize, pageOffset int) ([]*entities.Book, error)
	CountBooks(ctx context.Context, shelf string) (int, error)
}

func (l *ListBooksDomain) List(
	ctx context.Context,
	shelf string,
	pageSize, pageOffset int,
) (books []*entities.Book, finished bool, err error) {
	books, err = l.gateway.ListBooks(ctx, shelf, pageSize, pageOffset)
	if err != nil {
		return nil, false, fmt.Errorf("failed to list books in gateway: %w", err)
	}

	totalBooks, err := l.gateway.CountBooks(ctx, shelf)
	if err != nil {
		return nil, false, fmt.Errorf("failed to count books in gateway: %w", err)
	}

	finished = totalBooks <= pageOffset+pageSize

	return books, finished, nil
}
