package assistants

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

// AddAssistantLink добавляет ссылку ассистента в таблицу
func AddAssistantLink(DB *sql.DB, adminID uuid.UUID) (uuid.UUID, error) {
	// Генерация нового UUID для ассистента
	assistantID := uuid.New()

	// Генерация URL для чат-ассистента
	assistantURL := fmt.Sprintf("http://localhost/chat/%s", assistantID.String())

	// Вставляем запись в таблицу assistant_links
	query := `
        INSERT INTO assistant_links (assistant_id, id, url, created_at, updated_at)
        VALUES ($1, $2, $3, EXTRACT(epoch FROM now())::BIGINT, EXTRACT(epoch FROM now())::BIGINT)
        RETURNING id
    `
	var linkID uuid.UUID
	err := DB.QueryRow(query, adminID, assistantID, assistantURL).Scan(&linkID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to add assistant link: %w", err)
	}

	// Возвращаем UUID ссылки ассистента
	return linkID, nil
}
