package dtos

type BookCreateDto struct {
	Title    string `json:"title,omitempty" validate:"required,min=1,max=100"`
	Author   string `json:"author,omitempty" validate:"required,min=1,max=100"`
	Quantity int    `json:"quantity,omitempty" validate:"required"`
}

type BookUpdateDto struct {
	ID       string `json:"id,omitempty" validate:"required"`
	Title    string `json:"title,omitempty" validate:"required,min=1,max=100"`
	Author   string `json:"author,omitempty" validate:"required,min=1,max=100"`
	Quantity int    `json:"quantity,omitempty" validate:"required"`
}

type BookDto struct {
	ID       string `json:"id,omitempty"`
	Title    string `json:"title,omitempty"`
	Author   string `json:"author,omitempty"`
	Quantity int    `json:"quantity,omitempty"`
}
