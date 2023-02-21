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
	"github.com/sirupsen/logrus"
	echoSwagger "github.com/swaggo/echo-swagger"
	"os"
)

// @title           Echo Restful API
// @version         1.0
// @description     This is a sample restful server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api
func main() {
	e := echo.New()

	log := logrus.StandardLogger()

	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(file)
	log.SetLevel(logrus.InfoLevel)

	// to reach .env file
	_ = godotenv.Load()
	var env = os.Getenv("ENV")
	config := configs.GetConfig(env)
	mongoCollection := configs.ConnectDB(config.Database.Connection).
		Database(config.Database.DatabaseName).Collection(config.Database.CollectionName)

	// to create new repository with singleton pattern
	BookRepository := repository.GetSingleInstancesRepository(mongoCollection)

	// to create new service with singleton pattern
	BookService := service.GetSingleInstancesService(BookRepository)

	// to create new app
	app.NewBookHandler(e, BookService, log)

	// if we don't use this swagger give an error
	docs.SwaggerInfo.Host = "localhost:8080"
	// add swagger
	e.GET("/swagger/*any", echoSwagger.WrapHandler)

	// custom response
	//e.HTTPErrorHandler = app.NewHttpErrorHandler(models.NewErrorStatusCodeMaps()).Handler

	// start server
	log.Infof("Listening on port %s", config.Server.Port)
	e.Logger.Print(fmt.Sprintf("Listening on port %s", config.Server.Port))
	e.Logger.Fatal(e.Start(config.Server.Port))
}
