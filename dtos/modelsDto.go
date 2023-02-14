package dtos

type BookCreateRequest struct {
	Title    string `json:"title" validate:"required,min=1,max=100"`
	Author   string `json:"author" validate:"required,min=1,max=100"`
	Quantity int    `json:"quantity" validate:"required"`
}

type BookUpdateRequest struct {
	ID       string `json:"id" validate:"required"`
	Title    string `json:"title" validate:"required,min=1,max=100"`
	Author   string `json:"author" validate:"required,min=1,max=100"`
	Quantity int    `json:"quantity" validate:"required"`
}

type BookResponse struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Author   string `json:"author"`
	Quantity int    `json:"quantity"`
}

type BookCreateResponse struct {
	ID string `json:"id"`
}
