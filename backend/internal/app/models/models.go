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
	URL   string                  `form:"url" binding:"required"` // URL для ассистента
	Files []*multipart.FileHeader `form:"files[]"`                // Массив файлов
}

type SendMessageRequest struct {
	Message string `json:"message" binding:"required"`
}
