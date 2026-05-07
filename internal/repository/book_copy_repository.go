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

func (r *BookCopyRepository) GetByBook(ctx context.Context, bookID int) ([]models.BookCopyDetail, error) {
	query := `
		SELECT bc.id, bc.book_id, bc.branch_id, br.name, br.code, bc.status, bc.created_at
		FROM book_copies bc
		JOIN branches br ON br.id = bc.branch_id
		WHERE bc.book_id = $1
		ORDER BY br.name ASC, bc.id ASC
	`

	rows, err := r.DB.QueryContext(ctx, query, bookID)
	if err != nil {
		return nil, fmt.Errorf("query copies by book: %w", err)
	}
	defer rows.Close()

	var copies []models.BookCopyDetail
	for rows.Next() {
		var copy models.BookCopyDetail
		if err := rows.Scan(
			&copy.ID,
			&copy.BookID,
			&copy.BranchID,
			&copy.BranchName,
			&copy.BranchCode,
			&copy.Status,
			&copy.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan copy by book: %w", err)
		}
		copies = append(copies, copy)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate copies by book: %w", err)
	}

	return copies, nil
}
