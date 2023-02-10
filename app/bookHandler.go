package app

import (
	"RestfulWithEcho/dtos"
	"RestfulWithEcho/service"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"sync"
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

// with singleton pattern to create just one Handler we have to write like this or using once. Otherwise, every thread will create new Handler.
var lock = &sync.Mutex{}
var singleInstanceHandler *BookHandler

// for validation
var v = validator.New()

func GetSingleInstancesHandler(service service.IBookService) *BookHandler {
	if singleInstanceHandler == nil {
		lock.Lock()
		defer lock.Unlock()
		if singleInstanceHandler == nil {
			fmt.Println("Creating single handle instance now.")
			singleInstanceHandler = &BookHandler{Service: service}
		} else {
			fmt.Println("Single handle instance already created.")
		}
	} else {
		fmt.Println("Single handle instance already created.")
	}

	return singleInstanceHandler
}

// NewBookHandler => this method like constructor (#C) =>
/*func NewBookHandler(s service.IBookService) BookHandler {
	return BookHandler{Service: service}
}*/

// NewRouter => to create new route
func NewRouter(b *BookHandler) *echo.Echo {
	//Echo instance
	router := echo.New()

	router.Validator = &BookValidator{validator: v}

	//Routes
	router.GET("/books", b.GetAllBooks)
	router.GET("/books/:id", b.GetBookById)
	router.POST("/books", b.CreateBook)
	router.PUT("/books", b.UpdateBook)
	router.DELETE("/books/:id", b.DeleteBook)

	return router
}

// GetAllBooks => To get request for listing all of books

// GetAllBooks	 godoc
// @Summary      List books
// @Description  get all books
// @Tags         books
// @Accept       json
// @Produce      json
// @Param        q    query     string  false  "name search by q"  Format(email)
// @Router       /books [get]
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

	var bookDto dtos.BookCreateDto

	// We parse the data as json into the struct
	if err := c.Bind(&bookDto); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if err := c.Validate(bookDto); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	result, err := h.Service.Insert(bookDto)

	if err != nil || result == false {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, result)

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
