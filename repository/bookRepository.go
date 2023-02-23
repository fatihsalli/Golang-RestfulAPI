package repository

import (
	"RestfulWithEcho/models"
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type BookRepository struct {
	BookCollection *mongo.Collection
}

var singleInstanceRepo *BookRepository

func GetSingleInstancesRepository(mongoCollection *mongo.Collection) *BookRepository {
	if singleInstanceRepo == nil {
		fmt.Println("Creating single repository instance now.")
		singleInstanceRepo = &BookRepository{BookCollection: mongoCollection}
	} else {
		fmt.Println("Single repository instance already created.")
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

// Insert method => to create new book
func (b BookRepository) Insert(book models.Book) (bool, error) {
	// to open connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// mongodb.driver
	result, err := b.BookCollection.InsertOne(ctx, book)

	if result.InsertedID == nil || err != nil {
		return false, errors.New("failed to add")
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

	// => Update => update + insert = upsert => default value false
	// opt := options.Update().SetUpsert(true)
	filter := bson.D{{"_id", book.ID}}

	// => if we use this CreatedDate and id value will be null, so we have to use "UpdateOne"
	//replacement := models.Book{Title: book.Title, Quantity: book.Quantity, Author: book.Author, UpdatedDate: book.UpdatedDate}

	// => to update for one parameter
	//update := bson.D{{"$set", bson.D{{"title", book.Title}}}}

	// => if we have to chance more than one parameter we have to write like this
	update := bson.D{{"$set", bson.D{{"title", book.Title},
		{"author", book.Author}, {"quantity", book.Quantity}, {"updateddate", book.UpdatedDate}}}}

	// mongodb.driver
	result, err := b.BookCollection.UpdateOne(ctx, filter, update)

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

	//We can think of "Cursor" like a request. We pull the data from the database with the "Next" command. (C# => IQueryable)
	result, err := b.BookCollection.Find(ctx, bson.M{})

	if err != nil {
		return nil, err
	}

	for result.Next(ctx) {
		if err := result.Decode(&book); err != nil {
			return nil, err
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
	err := b.BookCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&book)

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
	result, err := b.BookCollection.DeleteOne(ctx, bson.M{"_id": id})

	if err != nil || result.DeletedCount <= 0 {
		return false, err
	}

	return true, nil

}
