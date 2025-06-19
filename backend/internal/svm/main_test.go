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

	// open file with name up
	stmt, err := os.ReadFile("C:/Users/achsanalfitra/Documents/hackathons/google-platform-hackathon/wandering/backend/internal/migration/0001_cannonical_order_up.sql")
	if err != nil {
		log.Fatalf("failed to load up statement")
	}

	if _, err := db.DB.Exec(string(stmt)); err != nil {
		log.Fatalf("failed to execute statement")
	}

	// run test code
	c := m.Run()

	// open file with name down
	stmt, err = os.ReadFile("C:/Users/achsanalfitra/Documents/hackathons/google-platform-hackathon/wandering/backend/internal/migration/0001_cannonical_order_down.sql")
	if err != nil {
		log.Fatalf("failed to load up statement")
	}

	if _, err := db.DB.Exec(string(stmt)); err != nil {
		log.Fatalf("failed to execute statement")
	}

	log.Println("success!")
	os.Exit(c)
}
