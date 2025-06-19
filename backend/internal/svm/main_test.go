package svm

import (
	"log"
	"os"
	"testing"

	"github.com/achsanalfitra/wandering/backend/internal/config"
	"github.com/achsanalfitra/wandering/backend/internal/util"
)

// testing entry point for the svm service

func TestMain(m *testing.M) {
	// load environment
	util.LoadEnv("C:/Users/achsanalfitra/Documents/hackathons/google-platform-hackathon/wandering/backend/internal/config/test/.env")

	// init database
	db := config.NewDatabase(
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("SSLMODE"),
	)

	err := db.DB.Ping()
	if err != nil {
		log.Fatalf("failed to ping once again")
	}

	log.Println("success!")
}
