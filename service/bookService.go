package service

import (
	"RestfulWithEcho/dtos"
	"RestfulWithEcho/models"
	"RestfulWithEcho/repository"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type BookService struct {
	Repository repository.IBookRepository
}

type IBookService interface {
	Insert(bookDto dtos.BookCreateDto) (bool, error)
	GetAll() ([]dtos.BookDto, error)
}

// NewBookService => to create new BookService
func NewBookService(repository repository.IBookRepository) BookService {
	return BookService{Repository: repository}
}

func (b BookService) Insert(bookDto dtos.BookCreateDto) (bool, error) {
	var book models.Book

	// we can use automapper, but it will cause performance loss.
	book.Title = bookDto.Title
	book.Quantity = bookDto.Quantity
	book.Author = bookDto.Author
	// to create id and created date value
	book.ID = uuid.New()
	book.CreatedDate = primitive.NewDateTimeFromTime(time.Now())

	result, err := b.Repository.Insert(book)

	if err != nil || result == false {
		return false, err
	}

	return true, nil
}

func (b BookService) GetAll() ([]dtos.BookDto, error) {
	result, err := b.Repository.GetAll()

	if err != nil {
		return nil, err
	}

	var bookDto dtos.BookDto
	var booksDto []dtos.BookDto

	for _, v := range result {
		bookDto.ID = v.ID
		bookDto.Title = v.Title
		bookDto.Quantity = v.Quantity
		bookDto.Author = v.Author

		booksDto = append(booksDto, bookDto)
	}

	return booksDto, nil
}
