package repository

import (
	"context"
	"database/sql"
	"fmt"

	"test1/internal/models"
)

type BookCopyRepository struct {
	DB *sql.DB
}

func NewBookCopyRepository(db *sql.DB) *BookCopyRepository {
	return &BookCopyRepository{DB: db}
}

func (r *BookCopyRepository) Create(ctx context.Context, req models.CreateBookCopyRequest) (*models.BookCopy, error) {
	query := `
		INSERT INTO book_copies (book_id, branch_id, status)
		VALUES ($1, $2, $3)
		RETURNING id, book_id, branch_id, status, created_at
	`

	var bc models.BookCopy
	err := r.DB.QueryRowContext(ctx, query, req.BookID, req.BranchID, req.Status).
		Scan(&bc.ID, &bc.BookID, &bc.BranchID, &bc.Status, &bc.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("insert book copy: %w", err)
	}

	return &bc, nil
}
