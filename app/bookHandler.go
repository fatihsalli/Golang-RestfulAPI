package app

import (
	"RestfulWithEcho/dtos"
	"RestfulWithEcho/service"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

// TODO:Swaggo kütüphanesi yüklenecek
type BookHandler struct {
	Service service.IBookService
}

// this method like constructor (#C) =>

func NewBookHandler(service service.IBookService) BookHandler {
	return BookHandler{Service: service}
}

// GetAllBooks => To get request for listing all of books
func (h BookHandler) GetAllBooks(c echo.Context) error {
	bookList, err := h.Service.GetAll()

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
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
			return c.JSON(http.StatusNotFound, err)
		}
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, book)

}

// CreateBook => To post request for creating new a book
func (h BookHandler) CreateBook(c echo.Context) error {

	var bookDto dtos.BookCreateDto

	// We parse the data as json into the struct
	if err := c.Bind(&bookDto); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	result, err := h.Service.Insert(bookDto)

	if err != nil || result == false {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, result)

}

// UpdateBook => To put request for changing exist book
func (h BookHandler) UpdateBook(c echo.Context) error {

	var bookDto dtos.BookUpdateDto

	// We parse the data as json into the struct
	if err := c.Bind(&bookDto); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	result, err := h.Service.Update(bookDto)

	if err != nil || result == false {
		return c.JSON(http.StatusInternalServerError, err)
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
		return c.JSON(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusOK, result)
}
