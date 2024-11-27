package main

import (
	"log"
	"t1/internal/pkg/db"
	"t1/internal/pkg/server"
	"os"
)

func main() {
	// Печатаем текущую рабочую директорию
	dir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get working directory: %v", err)
	}
	log.Printf("Current working directory: %s", dir)

	// Запуск сервера
	s := server.New(":8090")
	err = db.Migrate(s.DB, "/Users/chrizantona/t1-hack/backend/scripts/migration.sql")
	if err != nil {
		log.Fatalf("Migration failed: %v", err)
	}
	s.Start()
}
