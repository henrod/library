package entities

import "time"

type Book struct {
	ISBN       string
	Title      string
	Author     string
	CreateTime time.Time
	UpdateTime time.Time
}
