package main

import (
	"RestfulWithEcho/app"
	"RestfulWithEcho/configs"
	"RestfulWithEcho/docs"
	"RestfulWithEcho/repository"
	"RestfulWithEcho/service"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
	"os"
)

// @title Echo Restful API
// @description This is a sample restful server.

// @host localhost:8080
// @BasePath /api
func main() {
	// it has moved into app-book_handler.go
	e := echo.New()

	// to reach port => load env after that find port
	_ = godotenv.Load()
	//port := os.Getenv("PORT")

	var env = os.Getenv("ENV")

	config := configs.GetConfig(env)

	dbClient := configs.ConnectDB(config).Database("booksDB").Collection("books")

	// to create new repository with singleton pattern
	BookRepository := repository.GetSingleInstancesRepository(dbClient)

	// to create new service with singleton pattern
	BookService := service.GetSingleInstancesService(BookRepository)

	// to create new app with singleton pattern
	app.NewBookHandler(e, BookService)

	docs.SwaggerInfo.Host = "localhost:8080"
	// add swagger
	e.GET("/swagger/*any", echoSwagger.WrapHandler)
	// swag init -g main.go => bu komut denenecek

	// custom response
	//e.HTTPErrorHandler = app.NewHttpErrorHandler(models.NewErrorStatusCodeMaps()).Handler

	// start server
	e.Logger.Print(fmt.Sprintf("Listening on port %s", 8080))
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", 8080)))
}
