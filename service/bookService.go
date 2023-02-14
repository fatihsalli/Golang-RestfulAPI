package service

import (
	"RestfulWithEcho/dtos"
	"RestfulWithEcho/models"
	"RestfulWithEcho/repository"
	"fmt"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"sync"
	"time"
)

type BookService struct {
	Repository repository.IBookRepository
}

// Locka gerek yok zaten her request için single instance oluşturmalıyım
// with singleton pattern to create just one Service we have to write like this or using once. Otherwise, every thread will create new Service.
var lock = &sync.Mutex{}
var singleInstanceService *BookService

func GetSingleInstancesService(repository repository.IBookRepository) *BookService {
	if singleInstanceService == nil {
		lock.Lock()
		defer lock.Unlock()
		if singleInstanceService == nil {
			fmt.Println("Creating single service instance now.")
			singleInstanceService = &BookService{Repository: repository}
		} else {
			fmt.Println("Single service instance already created.")
		}
	} else {
		fmt.Println("Single service instance already created.")
	}

	return singleInstanceService
}

type IBookService interface {
	Insert(bookDto models.Book) (models.Book, error)
	GetAll() ([]models.Book, error)
	GetBookById(id string) (models.Book, error)
	Update(bookDto models.Book) (bool, error)
	Delete(id string) (bool, error)
}

func (b BookService) Insert(book models.Book) (models.Book, error) {

	// to create id and created date value
	book.ID = uuid.New().String()
	book.CreatedDate = primitive.NewDateTimeFromTime(time.Now())

	result, err := b.Repository.Insert(book)

	if err != nil || result == false {
		panic(err)
	}

	return book, nil
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

func (b BookService) GetBookById(id string) (dtos.BookDto, error) {
	var bookDto dtos.BookDto

	result, err := b.Repository.GetBookById(id)

	if err != nil {
		return bookDto, err
	}

	bookDto.ID = result.ID
	bookDto.Title = result.Title
	bookDto.Author = result.Author
	bookDto.Quantity = result.Quantity

	return bookDto, nil
}

func (b BookService) Update(bookDto dtos.BookUpdateDto) (bool, error) {
	var book models.Book

	// we can use automapper, but it will cause performance loss.
	book.ID = bookDto.ID
	book.Title = bookDto.Title
	book.Quantity = bookDto.Quantity
	book.Author = bookDto.Author
	// to create updated date value
	book.UpdatedDate = primitive.NewDateTimeFromTime(time.Now())

	result, err := b.Repository.Update(book)

	if err != nil || result == false {
		return false, err
	}

	return true, nil
}

func (b BookService) Delete(id string) (bool, error) {
	result, err := b.Repository.Delete(id)

	if err != nil || result == false {
		return false, err
	}

	return true, nil
}
