package repository

import (
	"RestfulWithEcho/models"
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"sync"
	"time"
)

type BookRepository struct {
	BookCollection *mongo.Collection
}

// with singleton pattern to create just one Repo we have to write like this or using once. Otherwise, every thread will create new Repo.
var lock = &sync.Mutex{}
var singleInstanceRepo *BookRepository

func GetSingleInstancesRepository(dbClient *mongo.Collection) *BookRepository {
	if singleInstanceRepo == nil {
		lock.Lock()
		defer lock.Unlock()
		if singleInstanceRepo == nil {
			fmt.Println("Creating single instance now.")
			singleInstanceRepo = &BookRepository{BookCollection: dbClient}
		} else {
			fmt.Println("Single instance already created.")
		}
	} else {
		fmt.Println("Single instance already created.")
	}

	return singleInstanceRepo
}

// IBookRepository to use for test or
type IBookRepository interface {
	Insert(book models.Book) (bool, error)
	GetAll() ([]models.Book, error)
	GetBookById(id string) (models.Book, error)
	Update(book models.Book) (bool, error)
	Delete(id string) (bool, error)
}

// NewBookRepository => this method like constructor (#C) =>
func NewBookRepository(dbClient *mongo.Collection) BookRepository {
	return BookRepository{BookCollection: dbClient}
}

// Insert method => to create new book
func (b BookRepository) Insert(book models.Book) (bool, error) {
	// to open connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// mongodb.driver
	result, err := b.BookCollection.InsertOne(ctx, book)

	if result.InsertedID == nil || err != nil {
		return false, errors.New("failed add")
	}

	return true, nil
}

// Update method => to change exist book
func (b BookRepository) Update(book models.Book) (bool, error) {
	// to open connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// to change updated date
	book.UpdatedDate = primitive.NewDateTimeFromTime(time.Now())

	// => Update => update + insert = upsert => opt is not necessary ???
	opt := options.Update().SetUpsert(true)
	filter := bson.D{{"id", book.ID}}

	// => if we use this CreatedDate and id value will be null, so we have to use "UpdateOne"
	//replacement := models.Book{Title: book.Title, Quantity: book.Quantity, Author: book.Author, UpdatedDate: book.UpdatedDate}

	// => to update for one parameter
	//update := bson.D{{"$set", bson.D{{"title", book.Title}}}}

	// => if we have to chance more than one parameter we have to write like this
	update := bson.D{{"$set", bson.D{{"title", book.Title},
		{"author", book.Author}, {"quantity", book.Quantity}, {"updateddate", book.UpdatedDate}}}}

	// mongodb.driver
	result, err := b.BookCollection.UpdateOne(ctx, filter, update, opt)

	if result.ModifiedCount <= 0 || err != nil {
		return false, errors.New("failed modify")
	}

	return true, nil
}

// GetAll Method => to list every books
func (b BookRepository) GetAll() ([]models.Book, error) {
	var book models.Book
	var books []models.Book

	// to open connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	//cursor olarak datayı çekerek Next ile database de gerekli datayı çekiyoruz(C# IQueryable gibi)
	result, err := b.BookCollection.Find(ctx, bson.M{})
	defer result.Close(ctx)

	if err != nil {
		log.Fatalln(err)
	}

	for result.Next(ctx) {
		if err := result.Decode(&book); err != nil {
			log.Fatalln(err)
		}
		// for appending book to books
		books = append(books, book)
	}

	return books, nil

}

// GetBookById Method => to find a single book with id
func (b BookRepository) GetBookById(id string) (models.Book, error) {
	var book models.Book

	// to open connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// to find book by id
	err := b.BookCollection.FindOne(ctx, bson.M{"id": id}).Decode(&book)

	if err != nil {
		return book, err
	}

	return book, nil
}

// Delete Method => to delete a book from books by id
func (b BookRepository) Delete(id string) (bool, error) {
	// to open connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// delete by id column
	result, err := b.BookCollection.DeleteOne(ctx, bson.M{"id": id})

	if err != nil || result.DeletedCount <= 0 {
		return false, err
	}

	return true, nil

}
