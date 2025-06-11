package repository

import (
	"JWTService/internal/models"
	"context"
	"database/sql"
	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO users (username, email, password_hash, created_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id`

	return r.db.QueryRowContext(ctx, query,
		user.Username,
		user.Email,
		user.Password_hash,
		user.CreatedAt,
	).Scan(&user.Id)
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	query := `SELECT * FROM users WHERE email = $1`

	err := r.db.GetContext(ctx, &user, query, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) IsEmailExists(ctx context.Context, email string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`

	err := r.db.GetContext(ctx, &exists, query, email)
	return exists, err
}
