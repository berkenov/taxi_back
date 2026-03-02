//go:build ignore
package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"
	}
	// Connect to default postgres database to create taxi
	dsn = "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Connect: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()
	_, err = db.Exec("CREATE DATABASE taxi")
	if err != nil {
		if err.Error() == `pq: database "taxi" already exists` {
			fmt.Println("Database taxi already exists")
			os.Exit(0)
		}
		fmt.Fprintf(os.Stderr, "Create DB: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Database taxi created")
}
