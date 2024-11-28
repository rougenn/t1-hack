package models

import (
	"mime/multipart"

	"github.com/google/uuid"
)

type Admin struct {
	ID           uuid.UUID `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	CreatedAt    int64     `json:"created_at"`
}

// Структура запроса для логина
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// Структура запроса для регистрации
type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

// Модель запроса для создания ассистента
type AssistantRequest struct {
	AssistantName      string                  `form:"assistant_name" binding:"required"` // имя для ассистента
	ModelName          string                  `form:"model_name" binding:"required"`     // имя для ассистента
	ChunkSize          string                  `form:"chunk_size" binding:"required"`
	EmbenddingsModelId string                  `form:"embeddings_model_id" binding:"required"`
	Files              []*multipart.FileHeader `form:"files[]"` // Массив файлов
}

type SendMessageRequest struct {
	Message string `json:"message" binding:"required"`
}

type ModelResponse struct {
	Message string `json:"message" binding:"required"`
}
