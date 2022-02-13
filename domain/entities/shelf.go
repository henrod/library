package entities

import "time"

type Shelf struct {
	Name       string
	CreateTime time.Time
	UpdateTime time.Time
}
