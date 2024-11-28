package db

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"

	_ "github.com/lib/pq"
)

func NewDB() *sql.DB {
	db, err := sql.Open("postgres", "postgres://username:password@localhost:5432/mydatabase?sslmode=disable")

	if err != nil {
		log.Fatal("connection: ", err)
		return nil
	}

	if err := db.Ping(); err != nil {
		log.Fatal("ping: ", err)
		return nil
	}

	return db
}

func Migrate(db *sql.DB, filePath string) error {
	// Чтение содержимого файла миграции
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read migration file: %v", err)
	}

	// Выполнение SQL миграции
	_, err = db.Exec(string(content))
	if err != nil {
		return fmt.Errorf("failed to execute migration: %v", err)
	}

	log.Println("Migration completed successfully")
	return nil
}
