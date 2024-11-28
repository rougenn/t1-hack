package main

import (
	"database/sql"
	"fmt"
	"log"
	"t1/internal/pkg/db"
	"t1/internal/pkg/server"
	"time"

	_ "github.com/lib/pq" // драйвер для PostgreSQL
)

// Функция для проверки доступности базы данных
func waitForDB(dsn string) {
	var dbConn *sql.DB
	var err error
	for i := 0; i < 10; i++ { // Пытаемся подключиться 10 раз
		dbConn, err = sql.Open("postgres", dsn)
		if err != nil {
			log.Println("Ошибка при подключении:", err)
			time.Sleep(5 * time.Second)
			continue
		}
		err = dbConn.Ping()
		if err == nil {
			fmt.Println("Подключение к базе данных успешно!")
			return
		}
		log.Println("Не удалось подключиться к базе данных:", err)
		time.Sleep(5 * time.Second)
	}
	log.Fatal("Не удалось подключиться к базе данных за 10 попыток")
}

func main() {
	// Формируем строку подключения к базе данных
	dsn := "postgres://user:password@db:5432/mydb?sslmode=disable"

	// Ожидаем подключения к базе данных
	waitForDB(dsn)

	// Формируем относительный путь к файлу миграции
	migrationFilePath := "scripts/migration.sql"

	// Запуск сервера
	s := server.New(":8090")

	// Запускаем миграцию
	err := db.Migrate(s.DB, migrationFilePath)
	if err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	// Запускаем сервер
	s.Start()
}
