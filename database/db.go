// database/db.go
package database

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// Connect establishes a connection to the database
func Connect(dbUser, dbName, dbPassword, dbHost string) (*sqlx.DB, error) {
	// Construct the connection string
	connStr := fmt.Sprintf("user=%s dbname=%s sslmode=disable password=%s host=%s",
		dbUser, dbName, dbPassword, dbHost)

	log.Printf("Attempting to connect to PostgreSQL: host=%s, user=%s, dbname=%s", dbHost, dbUser, dbName)

	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}
	
	log.Println("Successfully connected to database")
	return db, nil
}

// SetupTables creates necessary tables if they don't exist
func SetupTables(db *sqlx.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS todos (
			id SERIAL PRIMARY KEY,
			title TEXT NOT NULL,
			body TEXT,
			done BOOLEAN DEFAULT FALSE
		)
	`)

	if err != nil {
		log.Println("Error creating todos table:", err)
		return err
	}
	log.Println("Table 'todos' created or already exists")

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			username TEXT NOT NULL,
			email TEXT NOT NULL UNIQUE,
			password TEXT NOT NULL
		)
	`)

	if err != nil {
		log.Println("Error creating users table:", err)
		return err
	}
	log.Println("Table 'users' created or already exists")
	
	return nil
}