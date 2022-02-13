package entities

import "strings"

type Operation struct {
	Name       string
	Stage      string
	Percentage int
	Error      error
}

func (o *Operation) ResourceName() string {
	return strings.TrimPrefix(o.Name, "operations/")
}
