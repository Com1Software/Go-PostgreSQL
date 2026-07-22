package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	// Adjust these for your local PostgreSQL installation
	connStr := "host=localhost port=5432 user=postgres password=password dbname=yourdb sslmode=disable"

	db, err := sql.Open("pgx", connStr)
	if err != nil {
		log.Fatalf("Failed to open DB: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Cannot connect to DB: %v", err)
	}

	fmt.Println("Connected to PostgreSQL!")

	// 1. Create table
	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS users (
            id SERIAL PRIMARY KEY,
            name TEXT NOT NULL
        );
    `)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}
	fmt.Println("Table created.")

	// 2. Insert a record
	_, err = db.Exec(`INSERT INTO users (name) VALUES ($1);`, "Dave Example")
	if err != nil {
		log.Fatalf("Failed to insert record: %v", err)
	}
	fmt.Println("Record inserted.")

	// 3. Query and list records
	rows, err := db.Query(`SELECT id, name FROM users;`)
	if err != nil {
		log.Fatalf("Failed to query records: %v", err)
	}
	defer rows.Close()

	fmt.Println("Listing records:")
	for rows.Next() {
		var id int
		var name string
		if err := rows.Scan(&id, &name); err != nil {
			log.Fatalf("Scan error: %v", err)
		}
		fmt.Printf("ID=%d, Name=%s\n", id, name)
	}
}
