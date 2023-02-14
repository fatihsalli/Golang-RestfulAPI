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

	var router = e.Group("/books")
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
func (h BookHandler) GetAllBooks(c echo.Context) error {
	bookList, err := h.Service.GetAll()

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, bookList)
}

// GetBookById => To get request find a book by id
func (h BookHandler) GetBookById(c echo.Context) error {
	query := c.Param("id")

	// changed to objectId
	//cnv, _ := primitive.ObjectIDFromHex(query)

	book, err := h.Service.GetBookById(query)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.JSON(http.StatusNotFound, err.Error())
		}
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, book)

}

// CreateBook => To post request for creating new a book
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

	return c.JSON(http.StatusOK, bookCreateResponse)

}

// UpdateBook => To put request for changing exist book
func (h BookHandler) UpdateBook(c echo.Context) error {

	var bookDto dtos.BookUpdateDto

	// We parse the data as json into the struct
	if err := c.Bind(&bookDto); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	result, err := h.Service.Update(bookDto)

	if err != nil || result == false {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, result)

}

// DeleteBook => To delete request by id as a parameter
func (h BookHandler) DeleteBook(c echo.Context) error {
	query := c.Param("id")

	// changed to objectId
	//cnv, _ := primitive.ObjectIDFromHex(query)

	result, err := h.Service.Delete(query)

	if err != nil || result == false {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, result)
}
