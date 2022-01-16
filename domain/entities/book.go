package entities

import "time"

type Book struct {
	Name       string
	Author     string
	CreateTime time.Time
	UpdateTime time.Time
}
