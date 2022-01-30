package entities

import "time"

type Book struct {
	Name       string
	Author     string
	Shelf      *Shelf
	CreateTime time.Time
	UpdateTime time.Time
}
