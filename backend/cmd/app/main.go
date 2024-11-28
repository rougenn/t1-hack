package main

import (
	"log"
	"t1/internal/pkg/db"
	"t1/internal/pkg/server"
)

func main() {
	// Формируем относительный путь к файлу миграции
	migrationFilePath := "scripts/migration.sql"

	// Запуск сервера
	s := server.New(":8090")

	// Запускаем миграцию
	err := db.Migrate(s.DB, migrationFilePath)
	if err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	s.Start()
}
