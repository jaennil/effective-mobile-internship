package main

import (
	"database/sql"
	"fmt"
	"log"

    _ "github.com/lib/pq"
)

func main() {
	dbUser := "postgres"
	dbName := "postgres"
	dbPass := "password"
    host := "db"

	connectionString := fmt.Sprintf("password=%s user=%s dbname=%s sslmode=disable host=%s", dbPass, dbUser, dbName, host)
    fmt.Println(connectionString)

    db, err := sql.Open("postgres", connectionString)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    var result string
    err = db.QueryRow("SELECT 'Test'").Scan(&result)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("result!!!!!", result)
}
