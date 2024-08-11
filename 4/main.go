package main

import (
	"database/sql"
	"embed"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
    err := godotenv.Load()
    if err != nil {
	log.Fatal(err)
    }
    db := task1()
    defer db.Close()
    task9(db)
    task8(db)
    task2(db)
    task3(db, 1)
    task4(db, "DevOps")
    task5(db, "Analytics", 2)
    task6(db)
    task7()
}

func task1() *sql.DB {
    dbUrl := os.Getenv("DB_URL")
    if dbUrl == "" {
	log.Fatal("db url is empty")
    }

    db, err := sql.Open("postgres", dbUrl)
    if err != nil {
	log.Fatal(err)
    }

    err = db.Ping()
    if err != nil {
	log.Fatal(err)
    }

    log.Println("db ping success")
    return db
}

func task2(db *sql.DB) {
    rows, err := db.Query("SELECT name, last_name, department_id FROM employees")
    if err != nil {
	log.Fatal(err)
    }
    defer rows.Close()
    var (
	name, lastName string
	departmentId int
    )
    for rows.Next() {
	err := rows.Scan(&name, &lastName, &departmentId)
	if err != nil {
	    log.Fatal(err)
	}
	log.Printf("name: %s, last name: %s, department id: %d", name, lastName, departmentId)
    }
    err = rows.Err()
    if err != nil {
	log.Fatal(err)
    }
}

func task3(db *sql.DB, departmentId int) {
    rows, err := db.Query("SELECT name, last_name FROM employees WHERE department_id=$1", departmentId)
    if err != nil {
	log.Fatal(err)
    }
    defer rows.Close()
    var name, lastName string
    for rows.Next() {
	err := rows.Scan(&name, &lastName)
	if err != nil {
	    log.Fatal(err)
	}
	log.Printf("name: %s, last name: %s", name, lastName)
    }
    err = rows.Err()
    if err != nil {
	log.Fatal(err)
    }
}

func task4(db *sql.DB, departmentTitle string) {
    _, err := db.Exec("INSERT INTO departments(title) VALUES($1)", departmentTitle)
    if err != nil {
	log.Fatal(err)
    }
}

func task5(db *sql.DB, departmentTitle string, departmentId int) {
    _, err := db.Exec("UPDATE departments SET title = $1 WHERE department_id = $2", departmentTitle, departmentId)
    if err != nil {
	log.Fatal(err)
    }
}

func task6(db *sql.DB) {
    tx, err := db.Begin()
    if err != nil {
	log.Fatal(err)
    }
    _, err = tx.Exec("UPDATE departments SET title = 'Transaction test' WHERE department_id = 1")
    if err != nil {
	log.Fatal(err)
    }
    _, err = tx.Exec("INSERT INTO departments(title) VALUES('Transaction test2')")
    if err != nil {
	log.Fatal(err)
    }
    tx.Commit()
}

type Product struct {
	gorm.Model
	Code  string
	Price uint
}

func task7() {
    db, err := gorm.Open(postgres.Open(os.Getenv("DB_URL")), &gorm.Config{})
    if err != nil {
	    log.Fatal(err)
    }

    db.AutoMigrate(&Product{})

    db.Create(&Product{Code: "D42", Price: 100})

    var product Product
    db.First(&product, 4)
    db.First(&product, "code = ?", "D42")

    db.Model(&product).Update("Price", 200)

    db.Model(&product).Updates(Product{Price: 200, Code: "F42"})
    db.Model(&product).Updates(map[string]interface{}{"Price": 200, "Code": "F42"})

    db.Delete(&product, 1)
}

//go:embed migrations/*.sql
var embedMigrations embed.FS

func task8(db *sql.DB) {
    goose.SetBaseFS(embedMigrations)
    if err := goose.SetDialect("postgres"); err != nil {
	log.Fatal(err)
    }

    if err := goose.Up(db, "migrations"); err != nil {
	log.Fatal(err)
    }
}

func task9(db *sql.DB) {
    db.SetMaxOpenConns(5)
}
