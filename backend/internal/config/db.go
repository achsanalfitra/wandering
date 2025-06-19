package config

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type Database struct {
	port     string
	host     string
	sslmode  string
	user     string
	password string
	db       string
	DB       *sql.DB
}

func NewDatabase(user, password, db, port, host, sslmode string) *Database {
	database := Database{
		user:     user,
		password: password,
		db:       db,
		port:     port,
		host:     host,
		sslmode:  sslmode,
	}

	// check .env completeness
	if database.user == "" || database.password == "" || database.db == "" || database.port == "" || database.host == "" {
		log.Fatalf("missing one or more required PostgreSQL environment variables")
	}

	connStr := fmt.Sprintf("user=%s password=%s dbname=%s port=%s host=%s sslmode=%s",
		database.user, database.password, database.db, database.port, database.host, database.sslmode)

	dbConn, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("cannot establish connection to database")
	}

	// ping the db
	if err := dbConn.Ping(); err != nil {
		log.Fatalf("cannot ping the database")
	}

	// finally assign the connection
	database.DB = dbConn

	return &database
}
