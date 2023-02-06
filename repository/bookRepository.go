package repository

import (
	"RestfulWithEcho/models"
	"context"
	"errors"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"time"
)

type BookRepository struct {
	BookCollection *mongo.Collection
}

// IBookRepository to use for test or
type IBookRepository interface {
	Insert(book models.Book) (bool, error)
	GetAll() ([]models.Book, error)
	GetBookById(id primitive.ObjectID) (models.Book, error)
	Update(book models.Book) (bool, error)
	Delete(id primitive.ObjectID) (bool, error)
}

// Insert method => to create new book
func (b BookRepository) Insert(book models.Book) (bool, error) {
	// to open connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// to create id and created date value
	book.ID = uuid.New()
	book.CreatedDate = primitive.NewDateTimeFromTime(time.Now())

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

	filter := bson.D{{"id", book.ID}}

	// to change updated date
	book.UpdatedDate = primitive.NewDateTimeFromTime(time.Now())

	// mongodb.driver
	result, err := b.BookCollection.ReplaceOne(ctx, filter, book)

	if result.ModifiedCount <= 0 || err != nil {
		return false, errors.New("failed modify")
	}

	return true, nil
}

//// Update method => to change exist book
//func (b BookRepository) Update(book models.Book) (bool, error) {
//	// to open connection
//	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
//	defer cancel()
//
//	filter := bson.D{{"id", book.ID}}
//
//	// Update => update + insert = upsert olayını gerçekleştirebiliyoruz. Settings created datei değiştirmeden güncelleyebiliriz.
//	// mongodb.driver
//	result, err := b.BookCollection.InsertOne(ctx, filter, book)
//
//	if result.ModifiedCount <= 0 || err != nil {
//		return false, errors.New("failed modify")
//	}
//
//	return true, nil
//}

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
func (b BookRepository) GetBookById(id primitive.ObjectID) (models.Book, error) {
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
func (b BookRepository) Delete(id primitive.ObjectID) (bool, error) {
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

// this method like constructor (#C) =>

func NewBookRepository(dbClient *mongo.Collection) BookRepository {
	return BookRepository{BookCollection: dbClient}
}
