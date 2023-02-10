package dtos

type BookCreateDto struct {
	Title    string `json:"title" validate:"required,min=1,max=100"`
	Author   string `json:"author" validate:"required,min=1,max=100"`
	Quantity int    `json:"quantity" validate:"required"`
}

type BookUpdateDto struct {
	ID       string `json:"id" validate:"required"`
	Title    string `json:"title" validate:"required,min=1,max=100"`
	Author   string `json:"author" validate:"required,min=1,max=100"`
	Quantity int    `json:"quantity" validate:"required"`
}

type BookDto struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Author   string `json:"author"`
	Quantity int    `json:"quantity"`
}
