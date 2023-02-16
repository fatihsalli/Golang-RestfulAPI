package app

import (
	"RestfulWithEcho/dtos"
	"RestfulWithEcho/models"
	"RestfulWithEcho/service"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

type BookHandler struct {
	Service service.IBookService
}

// BookValidator echo validator for books
type BookValidator struct {
	validator *validator.Validate
}

// Validate validates books request body
func (b *BookValidator) Validate(i interface{}) error {
	return b.validator.Struct(i)
}

// for validation
var v = validator.New()

func NewBookHandler(e *echo.Echo, service service.IBookService) *BookHandler {

	router := e.Group("api/books")
	b := &BookHandler{Service: service}

	e.Validator = &BookValidator{validator: v}

	//Routes
	router.GET("", b.GetAllBooks)
	router.GET("/:id", b.GetBookById)
	router.POST("", b.CreateBook)
	router.PUT("", b.UpdateBook)
	router.DELETE("/:id", b.DeleteBook)

	return b
}

// GetAllBooks => To get request for listing all of books

// GetAllBooks godoc
// @Summary get all items in the book list
// @ID get-all-books
// @Produce json
// @Success 200 {array} dtos.BookResponse
// @Router /books [get]
func (h BookHandler) GetAllBooks(c echo.Context) error {
	bookList, err := h.Service.GetAll()

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	// we can use automapper, but it will cause performance loss.
	var bookResponse dtos.BookResponse
	var booksResponse []dtos.BookResponse

	for _, book := range bookList {
		bookResponse.ID = book.ID
		bookResponse.Title = book.Title
		bookResponse.Quantity = book.Quantity
		bookResponse.Author = book.Author
		booksResponse = append(booksResponse, bookResponse)
	}

	return c.JSON(http.StatusOK, booksResponse)
}

// GetBookById => To get request find a book by id

// GetBookById godoc
// @Summary get a book item by ID
// @ID get-book-by-id
// @Produce json
// @Param id path string true "book ID"
// @Success 200 {object} dtos.BookResponse
// @Router /books/{id} [get]
func (h BookHandler) GetBookById(c echo.Context) error {
	query := c.Param("id")

	book, err := h.Service.GetBookById(query)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.JSON(http.StatusNotFound, err.Error())
		}
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	// mapping
	var bookResponse dtos.BookResponse
	bookResponse.ID = book.ID
	bookResponse.Title = book.Title
	bookResponse.Author = book.Author
	bookResponse.Quantity = book.Quantity

	return c.JSON(http.StatusOK, bookResponse)

}

// CreateBook => To post request for creating new a book

// CreateBook godoc
// @Summary add a new item to the book list
// @ID create-book
// @Produce json
// @Param data body dtos.BookCreateRequest true "book data"
// @Success 201 {object} dtos.BookCreateResponse
// @Router /books [post]
func (h BookHandler) CreateBook(c echo.Context) error {

	var bookRequest dtos.BookCreateRequest

	// We parse the data as json into the struct
	if err := c.Bind(&bookRequest); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(bookRequest); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	var book models.Book

	// we can use automapper, but it will cause performance loss.
	book.Title = bookRequest.Title
	book.Quantity = bookRequest.Quantity
	book.Author = bookRequest.Author

	result, err := h.Service.Insert(book)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	// we have to return new id
	var bookCreateResponse dtos.BookCreateResponse
	bookCreateResponse.ID = result.ID

	return c.JSON(http.StatusCreated, bookCreateResponse)

}

// UpdateBook => To put request for changing exist book

// UpdateBook godoc
// @Summary update an item to the book list
// @ID update-book
// @Produce json
// @Param data body dtos.BookUpdateRequest true "book data"
// @Success 200 {object} bool
// @Router /books [put]
func (h BookHandler) UpdateBook(c echo.Context) error {

	var bookUpdateRequest dtos.BookUpdateRequest

	// We parse the data as json into the struct
	if err := c.Bind(&bookUpdateRequest); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	var book models.Book

	// we can use automapper, but it will cause performance loss.
	book.ID = bookUpdateRequest.ID
	book.Title = bookUpdateRequest.Title
	book.Quantity = bookUpdateRequest.Quantity
	book.Author = bookUpdateRequest.Author

	result, err := h.Service.Update(book)

	if err != nil || result == false {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, result)

}

// DeleteBook => To delete request by id as a parameter

// DeleteBook godoc
// @Summary delete a book item by ID
// @ID delete-book-by-id
// @Produce json
// @Param id path string true "book ID"
// @Success 200 {object} bool
// @Router /books/{id} [delete]
func (h BookHandler) DeleteBook(c echo.Context) error {
	query := c.Param("id")

	result, err := h.Service.Delete(query)

	if err != nil || result == false {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, result)
}
