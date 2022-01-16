package pg

import (
	"context"
	"fmt"
	"time"

	"github.com/go-pg/pg/v10"
)

type Gateway struct {
	db *pg.DB
}

func NewGateway(ctx context.Context, pgURL string) (*Gateway, error) {
	options, err := pg.ParseURL(pgURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse url '%s': %w", pgURL, err)
	}

	db := pg.Connect(options)

	timer := time.NewTimer(3 * time.Second)
	ticker := time.NewTicker(10 * time.Millisecond)
	defer timer.Stop()
	defer ticker.Stop()

waitConnection:
	for {
		select {
		case <-ticker.C:
			err = db.Ping(ctx)
			if err == nil {
				break waitConnection
			}
		case <-timer.C:
			return nil, fmt.Errorf("timeout connecting to postgres at %s", pgURL)
		}
	}

	return &Gateway{db: db}, nil
}
