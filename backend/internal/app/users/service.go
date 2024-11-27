package users

import (
	"database/sql"
	"errors"
	"fmt"
	"t1/internal/app/models"
	_ "time"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrIncorrectData = errors.New("password and/or username are incorrect")
	ErrAlreadyExists = errors.New("user with this email or phone already exists")
	ErrNotFound      = errors.New("user not found")
)

// SignIn аутентификация пользователя
func SignIn(DB *sql.DB, email, password string) (models.Admin, error) {
	user, err := GetUserByEmail(DB, email)
	if err != nil {
		return models.Admin{}, ErrIncorrectData
	}

	// Сравниваем хэшированный пароль
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return models.Admin{}, ErrIncorrectData
	}

	return user, nil
}

// Register регистрация пользователя
func Register(DB *sql.DB, req models.RegisterRequest) (models.Admin, error) {
	// Проверка на существование пользователя
	if _, err := GetUserByEmail(DB, req.Email); err == nil {
		return models.Admin{}, ErrAlreadyExists
	}

	// Генерация хеша пароля
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return models.Admin{}, fmt.Errorf("failed to hash password: %w", err)
	}

	// Создание нового пользователя
	user := models.Admin{
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
	}

	// Добавление пользователя в БД
	id, createdAt, err := AddToDB(DB, user)
	if err != nil {
		return models.Admin{}, err
	}

	// Присваиваем полученные значения ID и CreatedAt
	user.ID = id
	user.CreatedAt = createdAt
	return user, nil
}
