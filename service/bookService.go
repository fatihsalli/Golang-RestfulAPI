package service

import (
	"RestfulWithEcho/models"
	"RestfulWithEcho/repository"
	"fmt"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type BookService struct {
	Repository repository.IBookRepository
}

// with singleton pattern to create just one Service we have to write like this or using once. Otherwise, every thread will create new Service.
// var lock = &sync.Mutex{}
var singleInstanceService *BookService

func GetSingleInstancesService(repository repository.IBookRepository) *BookService {
	if singleInstanceService == nil {
		fmt.Println("Creating single service instance now.")
		singleInstanceService = &BookService{Repository: repository}
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
		return book, err
	}

	return book, nil
}

func (b BookService) GetAll() ([]models.Book, error) {
	result, err := b.Repository.GetAll()

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (b BookService) GetBookById(id string) (models.Book, error) {

	result, err := b.Repository.GetBookById(id)

	if err != nil {
		return result, err
	}

	return result, nil
}

func (b BookService) Update(book models.Book) (bool, error) {
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
