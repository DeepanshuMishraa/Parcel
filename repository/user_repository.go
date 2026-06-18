package repository

import (
	"context"
	"database/sql"
	"github.com/DeepanshuMishraa/mini-job-queue/models"
	"time"
)

func CreateUser(db *sql.DB, user models.User) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `INSERT INTO users(name,email,password) VALUES ($1,$2,$3) RETURNING id,email`

	created_users := &models.User{}

	err := db.QueryRowContext(ctx, query, user.Name, user.Email, user.Password).Scan(
		&created_users.Id,
		&created_users.Email,
	)

	if err != nil {
		return &models.User{}, err
	}

	return created_users, nil
}

func GetUserByEmail(db *sql.DB, email string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT id, name, email , password FROM users WHERE email = $1`

	users := &models.User{}

	err := db.QueryRowContext(ctx, query, email).Scan(
		&users.Id,
		&users.Name,
		&users.Email,
		&users.Password,
	)

	if err != nil {
		return nil, err
	}

	return users, nil
}
