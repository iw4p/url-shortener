package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Env struct {
	MongoURL   string
	DbName     string
	Collection string
}

func ReturnEnv() Env {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	return Env{
		MongoURL:   os.Getenv("MONGO_URL"),
		DbName:     os.Getenv("DB_NAME"),
		Collection: os.Getenv("COLLECTION_NAME"),
	}
}
