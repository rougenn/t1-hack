package models

import "time"

type Admin struct {
	ID           int       `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`          // хеш пароля
	CreatedAt    time.Time `json:"created_at"` // время создания как time.Time
}

type ChatRequest struct {
    UserID int    `json:"user_id" binding:"required"` // ID пользователя, который создает чат
    Title  string `json:"title" binding:"required"`   // Заголовок чата
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

// Модель чата
type Chat struct {
	ID        string    `json:"id"`        // ID чата (UUID)
	UserID    int       `json:"user_id"`    // ID пользователя, связанного с чатом
	Title     string    `json:"title"`      // Заголовок чата
	CreatedAt time.Time `json:"created_at"` // Дата и время создания чата
}

// Модель сообщения
type Message struct {
	ID        int       `json:"id"`        // ID сообщения
	ChatID    string    `json:"chat_id"`    // ID чата
	Sender    string    `json:"sender"`     // Отправитель сообщения
	Message   string    `json:"message"`    // Текст сообщения
	SentAt    time.Time `json:"sent_at"`    // Время отправки сообщения
}

// Структура для запроса на отправку сообщения
type SendMessageRequest struct {
	Message string `json:"message" binding:"required"`
}

// Структура для ответа о статусе отправки сообщения
type SendMessageResponse struct {
	ChatID        string    `json:"chat_id"`
	Message       string    `json:"message"`
	MessageContent string   `json:"message_content"`
}
