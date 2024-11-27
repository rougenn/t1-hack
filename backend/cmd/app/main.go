package main

import (
	"log"
	"os"
	"path/filepath"
	"t1/internal/pkg/db"
	"t1/internal/pkg/server"
)

func main() {
	// Печатаем текущую рабочую директорию
	dir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get working directory: %v", err)
	}
	log.Printf("Current working directory: %s", dir)

	// Формируем относительный путь к файлу миграции
	migrationFilePath := filepath.Join("scripts", "migration.sql")

	// Запуск сервера
	s := server.New(":8090")

	// Запускаем миграцию
	err = db.Migrate(s.DB, migrationFilePath)
	if err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	// Запуск сервера
	s.Start()
}
