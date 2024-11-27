package models

import (
	"mime/multipart"

	"github.com/gin-gonic/gin"
)

type Admin struct {
	ID           int    `json:"id"`
	Email        string `json:"email"`
	PasswordHash string `json:"-"`          // хеш пароля. будем сравнивать именно хеш.
	CreatedAt    int64  `json:"created_at"` // время создания юникс
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

// Модель запроса для создания ассистента
type AssistantRequest struct {
	URL   string          `json:"url" binding:"required"` // URL для ассистента
	Files *multipart.Form `json:"files"`                  // Файлы, передаваемые в запросе
	Ctx   *gin.Context    `json:"-"`                      // Контекст для доступа к файлам
}
