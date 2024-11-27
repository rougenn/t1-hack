package db

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/google/uuid"
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

// Функция для создания чата в базе данных
func CreateChatInDB(db *sql.DB, userID int, title string) (string, error) {
	// Проверяем, существует ли пользователь с таким ID
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)", userID).Scan(&exists)
	if err != nil || !exists {
		return "", fmt.Errorf("user with id %d does not exist", userID)
	}

	chatID := uuid.New().String() // Генерация уникального ID для чата

	query := `
        INSERT INTO chats (id, user_id, title)
        VALUES ($1, $2, $3)
        RETURNING id
    `

	err = db.QueryRow(query, chatID, userID, title).Scan(&chatID)
	if err != nil {
		log.Println("Error creating chat:", err)
		return "", err
	}

	return chatID, nil
}
