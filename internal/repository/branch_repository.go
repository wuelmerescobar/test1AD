package repository

import (
	"context"
	"database/sql"
	"fmt"

	"test1/internal/models"
)

type BranchRepository struct {
	DB *sql.DB
}

func NewBranchRepository(db *sql.DB) *BranchRepository {
	return &BranchRepository{DB: db}
}

func (r *BranchRepository) GetAll(ctx context.Context) ([]models.Branch, error) {
	query := `
		SELECT id, name, code, address, created_at
		FROM branches
		ORDER BY id ASC
	`

	rows, err := r.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query branches: %w", err)
	}
	defer rows.Close()

	var branches []models.Branch

	for rows.Next() {
		var b models.Branch
		if err := rows.Scan(&b.ID, &b.Name, &b.Code, &b.Address, &b.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan branch: %w", err)
		}
		branches = append(branches, b)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate branches: %w", err)
	}

	return branches, nil
}

func (r *BranchRepository) Create(ctx context.Context, req models.CreateBranchRequest) (*models.Branch, error) {
	query := `
		INSERT INTO branches (name, code, address)
		VALUES ($1, $2, $3)
		RETURNING id, name, code, address, created_at
	`

	var b models.Branch
	err := r.DB.QueryRowContext(ctx, query, req.Name, req.Code, req.Address).
		Scan(&b.ID, &b.Name, &b.Code, &b.Address, &b.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("insert branch: %w", err)
	}

	return &b, nil
}
