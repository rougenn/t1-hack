package assistants

import (
	"database/sql"
	"fmt"
)

// SaveMessageToDB сохраняет сообщение в базе данных
func SaveMessageToDB(db *sql.DB, assistantID, userID int, message string) error {
	query := `
		INSERT INTO message_history (assistant_id, user_id, message, created_at)
		VALUES ($1, $2, $3, NOW())
	`
	_, err := db.Exec(query, assistantID, userID, message)
	return err
}

// Функция для добавления ссылки ассистента в таблицу
func AddAssistantLink(DB *sql.DB, userID int, assistantID, url string) (int, error) {
	query := `
		INSERT INTO assistant_links (user_id, assistant_id, url)
		VALUES ($1, $2, $3)
		RETURNING id
	`
	var id int
	err := DB.QueryRow(query, userID, assistantID, url).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to add assistant link: %w", err)
	}
	return id, nil
}

// Функция для записи статистики запросов и ответов в таблицу
func AddAssistantStatistics(DB *sql.DB, linkID int, requestText, responseText string) error {
	query := `
		INSERT INTO assistant_statistics (link_id, request_text, response_text)
		VALUES ($1, $2, $3)
	`
	_, err := DB.Exec(query, linkID, requestText, responseText)
	if err != nil {
		return fmt.Errorf("failed to add assistant statistics: %w", err)
	}
	return nil
}
