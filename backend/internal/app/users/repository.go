package users

import (
	"database/sql"
	"t1/internal/app/models"
)

func AddToDB(db *sql.DB, user models.Admin) (int, int64, error) {
	query := `
    INSERT INTO admins (email, password_hash)
    VALUES ($1, $2)
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
	query := `DELETE FROM admins WHERE id = $1`
	_, err := db.Exec(query, userID)
	return err
}

func GetUserByID(db *sql.DB, userID int) (models.Admin, error) {
	var user models.Admin
	query := `
		SELECT id, email, password_hash, created_at
		FROM admins
		WHERE id = $1
	`
	row := db.QueryRow(query, userID)
	err := row.Scan(&user.ID, &user.Email, &user.PasswordHash, &user.CreatedAt)
	return user, err
}

func GetUserByEmail(db *sql.DB, email string) (models.Admin, error) {
	var user models.Admin
	query := `
		SELECT id, email, password_hash, created_at
		FROM admins
		WHERE email = $1
	`
	row := db.QueryRow(query, email)
	err := row.Scan(&user.ID, &user.Email, &user.PasswordHash, &user.CreatedAt)
	return user, err
}
