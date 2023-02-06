package main

import (
	"RestfulWithEcho/controller"
	"RestfulWithEcho/repository"
	"context"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

//TODO: #1: models klasörü yerine "models.go" ismi "models.go" olarak değiştirildi.
// ====> Soru klasör içinde olmayınca package ismi mainde kalıyor ve referans verirken hata alıyoruz???
// devamı ===> program ayağa kalkerken hata aldık

//TODO: #2: object id yerine uuid(v4) veya guid olarak değiştirildi.
//TODO: #3: createddate ve updateddate değerleri eklendi
//TODO: #4: controller tarafında da nesne üretme işi metot ile tanımlanacak (NewBookController)

func main() {
	e := echo.New()

	dbClient := ConnectDB().Database("booksDB").Collection("books")

	//TODO:Value yerine referans olarak verilecek
	BookRepository := repository.NewBookRepository(dbClient)

	// to create new controller
	BookController := controller.NewBookController(BookRepository)

	//TODO:controller constructor metotu içine // Routing group
	e.GET("/books", BookController.GetAllBooks)
	e.GET("/books/:id", BookController.GetBookById)
	e.POST("/books", BookController.CreateBook)
	e.PUT("/books", BookController.UpdateBook)
	e.DELETE("/books/:id", BookController.DeleteBook)
	e.Start(":8080")
}

// TODO: Ayrı bir yer ya da repository constructına
func ConnectDB() *mongo.Client {
	// we can use connection string => "mongodb://localhost:27017"
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))

	if err != nil {
		log.Fatalln(err)
	}
	// If don't connect within 20 seconds, give us an error
	var ctx, _ = context.WithTimeout(context.Background(), 20*time.Second)
	err = client.Connect(ctx)

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalln(err)
	}

	return client
}
