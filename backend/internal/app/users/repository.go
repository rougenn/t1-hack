package users

import (
	"database/sql"
	"t1/internal/app/models"
)

func AddToDB(db *sql.DB, user models.User) (int, int64, error) {
	query := `
    INSERT INTO users (first_name, second_name, company_name, email, phone_number, password_hash)
    VALUES ($1, $2, $3, $4, $5, $6)
    RETURNING id, created_at
	`
	var id int
	var createdAt int64
	err := db.QueryRow(query, user.Email, user.PasswordHash).Scan(&id, &createdAt)
	if err != nil {
		return 0, 0, err
	}
	return id, createdAt, nil
}

func DeleteFromDB(db *sql.DB, userID int) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := db.Exec(query, userID)
	return err
}

func GetUserByID(db *sql.DB, userID int) (models.User, error) {
	var user models.User
	query := `
		SELECT id, first_name, second_name, company_name, email, phone_number, password_hash, created_at
		FROM users
		WHERE id = $1
	`
	row := db.QueryRow(query, userID)
	err := row.Scan(&user.ID, &user.Email, &user.PasswordHash, &user.CreatedAt)
	return user, err
}

func GetUserByPhone(db *sql.DB, phone string) (models.User, error) {
	var user models.User
	query := `
		SELECT id, first_name, second_name, company_name, email, phone_number, password_hash, created_at
		FROM users
		WHERE phone_number = $1
	`
	row := db.QueryRow(query, phone)
	err := row.Scan(&user.ID, &user.Email, &user.PasswordHash, &user.CreatedAt)
	return user, err
}

func GetUserByEmail(db *sql.DB, email string) (models.User, error) {
	var user models.User
	query := `
		SELECT id, first_name, second_name, company_name, email, phone_number, password_hash, created_at
		FROM users
		WHERE email = $1
	`
	row := db.QueryRow(query, email)
	err := row.Scan(&user.ID, &user.Email, &user.PasswordHash, &user.CreatedAt)
	return user, err
}
