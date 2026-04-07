package repository

import (
	"context"
	"database/sql"
	"fmt"

	"test1/internal/models"
)

type AccountRepository struct {
	DB *sql.DB
}

func NewAccountRepository(db *sql.DB) *AccountRepository {
	return &AccountRepository{DB: db}
}

func (r *AccountRepository) Create(ctx context.Context, email, passwordHash, role string) (*models.Account, error) {
	query := `
		INSERT INTO accounts (email, password_hash, role)
		VALUES ($1, $2, $3)
		RETURNING id, email, password_hash, role, is_active, created_at
	`

	var a models.Account
	err := r.DB.QueryRowContext(ctx, query, email, passwordHash, role).
		Scan(&a.ID, &a.Email, &a.PasswordHash, &a.Role, &a.IsActive, &a.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("create account: %w", err)
	}

	return &a, nil
}

func (r *AccountRepository) GetByEmail(ctx context.Context, email string) (*models.Account, error) {
	query := `
		SELECT id, email, password_hash, role, is_active, created_at
		FROM accounts
		WHERE email = $1
	`

	var a models.Account
	err := r.DB.QueryRowContext(ctx, query, email).
		Scan(&a.ID, &a.Email, &a.PasswordHash, &a.Role, &a.IsActive, &a.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("account not found")
		}
		return nil, fmt.Errorf("get account by email: %w", err)
	}

	return &a, nil
}
