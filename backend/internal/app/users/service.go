package users

import (
	"database/sql"
	"errors"
	"fmt"
	"t1/internal/app/models"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrIncorrectData = errors.New("password and/or username are incorrect")
	ErrAlreadyExists = errors.New("user with this email or phone already exists")
	ErrNotFound      = errors.New("user not found")
)

// SignIn аутентификация пользователя
func SignIn(DB *sql.DB, email, password string) (models.User, error) {
	user, err := GetUserByEmail(DB, email)
	if err != nil {
		return models.User{}, ErrIncorrectData
	}

	// Сравниваем хэшированный пароль
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return models.User{}, ErrIncorrectData
	}

	return user, nil
}

// Register регистрация пользователя
func Register(DB *sql.DB, req models.RegisterRequest) (models.User, error) {
	if _, err := GetUserByEmail(DB, req.Email); err == nil {
		return models.User{}, ErrAlreadyExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return models.User{}, fmt.Errorf("failed to hash password: %w", err)
	}

	// Создаем пользователя
	user := models.User{
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
	}

	id, createdAt, err := AddToDB(DB, user)
	if err != nil {
		return models.User{}, err
	}

	user.ID = id
	user.CreatedAt = createdAt
	return user, nil
}
