package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
)

// DialPostgresDB is a Wire provider function that connects to a Postgres database
func DialPostgresDB() (*sql.DB, func(), error) {
	host := os.Getenv("PG_HOST")
	if host == "" {
		host = "localhost"
	}

	database := os.Getenv("PG_DATABASE")
	if database == "" {
		database = "oauthfarm"
	}

	user := os.Getenv("PG_USER")
	if user == "" {
		user = "tanopwan"
	}

	password := os.Getenv("PG_PASSWORD")

	fmt.Println("[DialPostgresDB] connecting to postgresql ", host)

	connStr := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", user, password, host, database)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, nil, err
	}

	cleanup := func() {
		fmt.Print("[DialPostgresDB] Closing connection")
		if err := db.Close(); err != nil {
			log.Fatalf("error closing db connection with reason: %s\n", err.Error())
		}
	}

	return db, cleanup, nil
}
