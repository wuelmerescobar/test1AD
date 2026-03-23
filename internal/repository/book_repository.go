package repository

import (
	"context"
	"database/sql"
	"fmt"

	"test1/internal/models"
)

type BookRepository struct {
	DB *sql.DB
}

func NewBookRepository(db *sql.DB) *BookRepository {
	return &BookRepository{DB: db}
}

func (r *BookRepository) GetAll(ctx context.Context) ([]models.Book, error) {
	query := `
		SELECT id, title, author, isbn, genre, created_at
		FROM books
		ORDER BY id ASC
	`

	rows, err := r.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query books: %w", err)
	}
	defer rows.Close()

	var books []models.Book

	for rows.Next() {
		var b models.Book
		if err := rows.Scan(&b.ID, &b.Title, &b.Author, &b.ISBN, &b.Genre, &b.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan book: %w", err)
		}
		books = append(books, b)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate books: %w", err)
	}

	return books, nil
}

func (r *BookRepository) Create(ctx context.Context, req models.CreateBookRequest) (*models.Book, error) {
	query := `
		INSERT INTO books (title, author, isbn, genre)
		VALUES ($1, $2, $3, $4)
		RETURNING id, title, author, isbn, genre, created_at
	`

	var b models.Book
	err := r.DB.QueryRowContext(ctx, query, req.Title, req.Author, req.ISBN, req.Genre).
		Scan(&b.ID, &b.Title, &b.Author, &b.ISBN, &b.Genre, &b.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("insert book: %w", err)
	}

	return &b, nil
}

func (r *BookRepository) GetAvailableByBranch(ctx context.Context, branchID int) ([]models.Book, error) {
	query := `
		SELECT b.id, b.title, b.author, b.isbn, b.genre, b.created_at
		FROM books b
		JOIN book_copies bc ON b.id = bc.book_id
		WHERE bc.branch_id = $1 AND bc.status = 'available'
		ORDER BY b.id ASC
	`

	rows, err := r.DB.QueryContext(ctx, query, branchID)
	if err != nil {
		return nil, fmt.Errorf("query available books: %w", err)
	}
	defer rows.Close()

	var books []models.Book

	for rows.Next() {
		var b models.Book
		if err := rows.Scan(&b.ID, &b.Title, &b.Author, &b.ISBN, &b.Genre, &b.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan book: %w", err)
		}
		books = append(books, b)
	}

	return books, nil
}

func (r *BookRepository) Update(ctx context.Context, id int, req models.CreateBookRequest) (*models.Book, error) {
	query := `
		UPDATE books
		SET title = $1, author = $2, isbn = $3, genre = $4
		WHERE id = $5
		RETURNING id, title, author, isbn, genre, created_at
	`

	var b models.Book
	err := r.DB.QueryRowContext(ctx, query, req.Title, req.Author, req.ISBN, req.Genre, id).
		Scan(&b.ID, &b.Title, &b.Author, &b.ISBN, &b.Genre, &b.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("book not found")
		}
		return nil, fmt.Errorf("update book: %w", err)
	}

	return &b, nil
}

func (r *BookRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM books WHERE id = $1`

	result, err := r.DB.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete book: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("book not found")
	}

	return nil
}
