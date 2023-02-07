package dtos

import (
	"github.com/google/uuid"
)

type BookCreateDto struct {
	Title    string `json:"title,omitempty"`
	Author   string `json:"author,omitempty"`
	Quantity int    `json:"quantity,omitempty"`
}

type BookUpdateDto struct {
	ID       uuid.UUID `json:"id,omitempty"`
	Title    string    `json:"title,omitempty"`
	Author   string    `json:"author,omitempty"`
	Quantity int       `json:"quantity,omitempty"`
}

type BookDto struct {
	ID       uuid.UUID `json:"id,omitempty"`
	Title    string    `json:"title,omitempty"`
	Author   string    `json:"author,omitempty"`
	Quantity int       `json:"quantity,omitempty"`
}
