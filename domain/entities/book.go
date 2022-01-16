package entities

import "time"

type Book struct {
	Title      string
	Author     string
	CreateTime time.Time
	UpdateTime time.Time
}
