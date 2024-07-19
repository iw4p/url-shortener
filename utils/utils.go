package utils

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type EnvironmentVariables struct {
	Mongo_url  string
	Db_name    string
	Collection string
}

func ReturnEnv() EnvironmentVariables {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	env := EnvironmentVariables{
		Mongo_url:  os.Getenv("MONGODB_URI"),
		Db_name:    os.Getenv("DB_NAME"),
		Collection: os.Getenv("COLLECTION"),
	}
	return env
}
