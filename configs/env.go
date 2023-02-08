package configs

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

func EnvMongoURI() string {
	err := godotenv.Load()

	if err != nil {
		log.Fatalln("error .env")
	}

	// with "go-dotenv" package we can read inside ".env" => mongodb://localhost:27017
	mongoIRU := os.Getenv("MONGOURI")
	return mongoIRU
}
