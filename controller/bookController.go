package controller

import (
	"RestfulWithEcho/models"
	"RestfulWithEcho/repository"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

// TODO:Service layer
// TODO:Swaggo kütüphanesi yüklenecek
type Controller struct {
	Repo repository.IBookRepository
}

// this method like constructor (#C) =>

func NewBookController(repo repository.IBookRepository) Controller {
	return Controller{Repo: repo}
}

// GetAllBooks => To get request for listing all of books
func (cont Controller) GetAllBooks(c echo.Context) error {
	bookList, err := cont.Repo.GetAll()

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, bookList)
}

// GetBookById => To get request find a book by id
func (cont Controller) GetBookById(c echo.Context) error {
	query := c.Param("id")

	// changed to objectId
	cnv, _ := primitive.ObjectIDFromHex(query)

	book, err := cont.Repo.GetBookById(cnv)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.JSON(http.StatusNotFound, err)
		}
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, book)

}

// CreateBook => To post request for creating new a book
func (cont Controller) CreateBook(c echo.Context) error {

	var book models.Book

	// We parse the data as json into the struct
	if err := c.Bind(&book); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	result, err := cont.Repo.Insert(book)

	if err != nil || result == false {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, result)

}

// UpdateBook => To put request for changing exist book
func (cont Controller) UpdateBook(c echo.Context) error {

	var book models.Book

	// We parse the data as json into the struct
	if err := c.Bind(&book); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	result, err := cont.Repo.Update(book)

	if err != nil || result == false {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, result)

}

// DeleteBook => To delete request by id as a parameter
func (cont Controller) DeleteBook(c echo.Context) error {
	query := c.Param("id")

	// changed to objectId
	cnv, _ := primitive.ObjectIDFromHex(query)

	result, err := cont.Repo.Delete(cnv)

	if err != nil || result == false {
		return c.JSON(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusOK, result)
}
