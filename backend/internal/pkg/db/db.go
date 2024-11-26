package db

import (
	"database/sql"
	"io/ioutil"
	"log"

	_ "github.com/lib/pq"
)

func NewDB() *sql.DB {
	// get username, password and bdname from environment
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
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	_, err = db.Exec(string(content))
	if err != nil {
		return err
	}

	log.Println("Migration completed successfully")
	return nil
}
