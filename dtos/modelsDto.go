package dtos

type BookCreateDto struct {
	Title    string `json:"title,omitempty"`
	Author   string `json:"author,omitempty"`
	Quantity int    `json:"quantity,omitempty"`
}

type BookUpdateDto struct {
	ID       string `json:"id,omitempty"`
	Title    string `json:"title,omitempty"`
	Author   string `json:"author,omitempty"`
	Quantity int    `json:"quantity,omitempty"`
}

type BookDto struct {
	ID       string `json:"id,omitempty"`
	Title    string `json:"title,omitempty"`
	Author   string `json:"author,omitempty"`
	Quantity int    `json:"quantity,omitempty"`
}
