package books

import (
	"context"
	"fmt"

	"github.com/Henrod/library/domain/errors"

	"google.golang.org/protobuf/types/known/fieldmaskpb"

	"github.com/Henrod/library/domain/entities"
)

type UpdateBookDomain struct {
	gateway UpdateBookGateway
}

var notUserUpdatableFields = map[string]struct{}{
	"create_time": {},
	"update_time": {},
}

func NewUpdateBookDomain(gateway UpdateBookGateway) *UpdateBookDomain {
	return &UpdateBookDomain{gateway: gateway}
}

type UpdateBookGateway interface {
	UpdateBook(ctx context.Context, shelfName string, book *entities.Book, fields []string) (*entities.Book, error)
}

func (g *UpdateBookDomain) UpdateBook(
	ctx context.Context,
	shelfName string,
	inputBook *entities.Book,
	updateMask *fieldmaskpb.FieldMask,
) (*entities.Book, error) {
	if updateMask == nil {
		return nil, &errors.BadRequestError{
			InvalidField: "update_mask",
			Details:      "update_mask must contain book fields (<link to swagger book resource>)",
		}
	}

	updateMask.Normalize()

	fields := make([]string, 0)
	for _, path := range updateMask.GetPaths() {
		if _, ok := notUserUpdatableFields[path]; ok {
			continue
		}

		fields = append(fields, path)
	}

	if len(fields) == 0 {
		return nil, &errors.BadRequestError{
			InvalidField: "update_mask",
			Details:      "update_mask doesn't have any valid fields to update (<link to swagger book resource>)",
		}
	}

	book, err := g.gateway.UpdateBook(ctx, shelfName, inputBook, fields)
	if err != nil {
		return nil, fmt.Errorf("failed to update book in gateway: %w", err)
	}

	if book == nil {
		return nil, errors.NotFoundError{
			Details: fmt.Sprintf("book %s at shelf %s not found", shelfName, inputBook.Name),
		}
	}

	return book, nil
}
