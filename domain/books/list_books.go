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
	ListBooks(ctx context.Context, shelf string) ([]*entities.Book, error)
}

func (l *ListBooksDomain) List(ctx context.Context, shelf string) ([]*entities.Book, error) {
	books, err := l.gateway.ListBooks(ctx, shelf)
	if err != nil {
		return nil, fmt.Errorf("failed to list books in gateway: %w", err)
	}

	return books, nil
}
