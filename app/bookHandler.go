package app

import (
	"RestfulWithEcho/dtos"
	"RestfulWithEcho/errors"
	"RestfulWithEcho/models"
	"RestfulWithEcho/response"
	"RestfulWithEcho/service"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

type BookHandler struct {
	Service service.IBookService
	Logger  *logrus.Logger
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

func NewBookHandler(e *echo.Echo, service service.IBookService, logger *logrus.Logger) *BookHandler {

	router := e.Group("api/books")
	b := &BookHandler{Service: service, Logger: logger}

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
// @Success 200 {array} response.JSONSuccessResultData
// @Success 500 {object} errors.InternalServerError
// @Router /books [get]
func (h BookHandler) GetAllBooks(c echo.Context) error {
	bookList, err := h.Service.GetAll()

	if err != nil {
		h.Logger.Errorf("StatusInternalServerError: %v", err.Error())
		return c.JSON(http.StatusInternalServerError, errors.InternalServerError{
			Message: "Something went wrong!",
		})
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

	// to response success result data
	jsonSuccessResultData := response.JSONSuccessResultData{
		TotalItemCount: len(booksResponse),
		Data:           booksResponse,
	}

	return c.JSON(http.StatusOK, jsonSuccessResultData)
}

// GetBookById => To get request find a book by id

// GetBookById godoc
// @Summary get a book item by ID
// @ID get-book-by-id
// @Produce json
// @Param id path string true "book ID"
// @Success 200 {object} response.JSONSuccessResultData
// @Success 404 {object} errors.NotFoundError
// @Success 500 {object} errors.InternalServerError
// @Router /books/{id} [get]
func (h BookHandler) GetBookById(c echo.Context) error {
	query := c.Param("id")

	book, err := h.Service.GetBookById(query)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			h.Logger.Errorf("Not found exception: {%v} with id not found!", query)
			return c.JSON(http.StatusNotFound, errors.NotFoundError{
				Message: fmt.Sprintf("Not found exception: {%v} with id not found!", query),
			})
		}
		h.Logger.Errorf("StatusInternalServerError: %v", err.Error())
		return c.JSON(http.StatusInternalServerError, errors.InternalServerError{
			Message: "Something went wrong!",
		})
	}

	// mapping
	var bookResponse dtos.BookResponse
	bookResponse.ID = book.ID
	bookResponse.Title = book.Title
	bookResponse.Author = book.Author
	bookResponse.Quantity = book.Quantity

	// to response success result data => single one
	jsonSuccessResultData := response.JSONSuccessResultData{
		TotalItemCount: 1,
		Data:           bookResponse,
	}

	return c.JSON(http.StatusOK, jsonSuccessResultData)

}

// CreateBook => To post request for creating new a book

// CreateBook godoc
// @Summary add a new item to the book list
// @ID create-book
// @Produce json
// @Param data body dtos.BookCreateRequest true "book data"
// @Success 201 {object} response.JSONSuccessResultId
// @Success 400 {object} errors.BadRequestError
// @Success 500 {object} errors.InternalServerError
// @Router /books [post]
func (h BookHandler) CreateBook(c echo.Context) error {

	var bookRequest dtos.BookCreateRequest

	// We parse the data as json into the struct
	if err := c.Bind(&bookRequest); err != nil {
		h.Logger.Errorf("Bad Request. It cannot be binding! %v", err.Error())
		return c.JSON(http.StatusBadRequest, errors.BadRequestError{
			Message: fmt.Sprintf("Bad Request. It cannot be binding! %v", err.Error()),
		})
	}

	if err := c.Validate(bookRequest); err != nil {
		h.Logger.Errorf("Bad Request! %v", err.Error())
		return c.JSON(http.StatusBadRequest, errors.BadRequestError{
			Message: fmt.Sprintf("Bad Request! %v", err.Error()),
		})
	}

	var book models.Book

	// we can use automapper, but it will cause performance loss.
	book.Title = bookRequest.Title
	book.Quantity = bookRequest.Quantity
	book.Author = bookRequest.Author

	result, err := h.Service.Insert(book)

	if err != nil {
		h.Logger.Errorf("StatusInternalServerError: %v", err.Error())
		return c.JSON(http.StatusInternalServerError, &errors.InternalServerError{
			Message: "Book cannot create! Something went wrong.",
		})
	}

	// to response id and success boolean
	jsonSuccessResultId := response.JSONSuccessResultId{
		ID:      result.ID,
		Success: true,
	}

	return c.JSON(http.StatusCreated, jsonSuccessResultId)

}

// UpdateBook => To put request for changing exist book

// UpdateBook godoc
// @Summary update an item to the book list
// @ID update-book
// @Produce json
// @Param data body dtos.BookUpdateRequest true "book data"
// @Success 200 {object} response.JSONSuccessResultId
// @Success 400 {object} errors.BadRequestError
// @Success 500 {object} errors.InternalServerError
// @Router /books [put]
func (h BookHandler) UpdateBook(c echo.Context) error {

	var bookUpdateRequest dtos.BookUpdateRequest

	// we parse the data as json into the struct
	if err := c.Bind(&bookUpdateRequest); err != nil {
		h.Logger.Errorf("Bad Request! %v", err)
		return c.JSON(http.StatusBadRequest, errors.BadRequestError{
			Message: fmt.Sprintf("Bad Request. It cannot be binding! %v", err.Error()),
		})
	}

	// validation
	if err := c.Validate(bookUpdateRequest); err != nil {
		h.Logger.Errorf("Bad Request! %v", err)
		return c.JSON(http.StatusBadRequest, errors.BadRequestError{
			Message: fmt.Sprintf("Bad Request! %v", err.Error()),
		})
	}

	if _, err := h.Service.GetBookById(bookUpdateRequest.ID); err != nil {
		h.Logger.Errorf("Not found exception: {%v} with id not found!", bookUpdateRequest.ID)
		return c.JSON(http.StatusNotFound, errors.NotFoundError{
			Message: fmt.Sprintf("Not found exception: {%v} with id not found!", bookUpdateRequest.ID),
		})
	}

	var book models.Book

	// we can use automapper, but it will cause performance loss.
	book.ID = bookUpdateRequest.ID
	book.Title = bookUpdateRequest.Title
	book.Quantity = bookUpdateRequest.Quantity
	book.Author = bookUpdateRequest.Author

	result, err := h.Service.Update(book)

	if err != nil || result == false {
		h.Logger.Errorf("StatusInternalServerError: {%v} ", err.Error())
		return c.JSON(http.StatusInternalServerError, &errors.InternalServerError{
			Message: "Book cannot create! Something went wrong.",
		})
	}

	// to response id and success boolean
	jsonSuccessResultId := response.JSONSuccessResultId{
		ID:      book.ID,
		Success: result,
	}

	return c.JSON(http.StatusOK, jsonSuccessResultId)

}

// DeleteBook => To delete request by id as a parameter

// DeleteBook godoc
// @Summary delete a book item by ID
// @ID delete-book-by-id
// @Produce json
// @Param id path string true "book ID"
// @Success 200 {object} response.JSONSuccessResultId
// @Success 404 {object} errors.NotFoundError
// @Router /books/{id} [delete]
func (h BookHandler) DeleteBook(c echo.Context) error {
	query := c.Param("id")

	result, err := h.Service.Delete(query)

	if err != nil || result == false {
		h.Logger.Errorf("Not found exception: {%v} with id not found!", query)
		return c.JSON(http.StatusNotFound, errors.NotFoundError{
			Message: fmt.Sprintf("Not found exception: {%v} with id not found!", query),
		})
	}

	// to response id and success boolean
	jsonSuccessResultId := response.JSONSuccessResultId{
		ID:      query,
		Success: result,
	}

	return c.JSON(http.StatusOK, jsonSuccessResultId)
}
