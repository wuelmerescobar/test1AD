package repository

import (
	"context"
	"database/sql"
	"fmt"

	"test1/internal/models"
)

type MemberRepository struct {
	DB *sql.DB
}

func NewMemberRepository(db *sql.DB) *MemberRepository {
	return &MemberRepository{DB: db}
}

func (r *MemberRepository) GetAll(ctx context.Context) ([]models.Member, error) {
	query := `
		SELECT id, first_name, last_name, email, phone, branch_id, created_at
		FROM members
		ORDER BY id ASC
	`

	rows, err := r.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query members: %w", err)
	}
	defer rows.Close()

	var members []models.Member

	for rows.Next() {
		var m models.Member
		if err := rows.Scan(
			&m.ID,
			&m.FirstName,
			&m.LastName,
			&m.Email,
			&m.Phone,
			&m.BranchID,
			&m.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan member: %w", err)
		}
		members = append(members, m)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate members: %w", err)
	}

	return members, nil
}

func (r *MemberRepository) GetByBranch(ctx context.Context, branchID int) ([]models.Member, error) {
	query := `
		SELECT id, first_name, last_name, email, phone, branch_id, created_at
		FROM members
		WHERE branch_id = $1
		ORDER BY id ASC
	`

	rows, err := r.DB.QueryContext(ctx, query, branchID)
	if err != nil {
		return nil, fmt.Errorf("query members by branch: %w", err)
	}
	defer rows.Close()

	var members []models.Member

	for rows.Next() {
		var m models.Member
		if err := rows.Scan(
			&m.ID,
			&m.FirstName,
			&m.LastName,
			&m.Email,
			&m.Phone,
			&m.BranchID,
			&m.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan member by branch: %w", err)
		}
		members = append(members, m)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate members by branch: %w", err)
	}

	return members, nil
}

func (r *MemberRepository) Create(ctx context.Context, req models.CreateMemberRequest) (*models.Member, error) {
	query := `
		INSERT INTO members (first_name, last_name, email, phone, branch_id)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, first_name, last_name, email, phone, branch_id, created_at
	`

	var m models.Member
	err := r.DB.QueryRowContext(
		ctx,
		query,
		req.FirstName,
		req.LastName,
		req.Email,
		req.Phone,
		req.BranchID,
	).Scan(
		&m.ID,
		&m.FirstName,
		&m.LastName,
		&m.Email,
		&m.Phone,
		&m.BranchID,
		&m.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("insert member: %w", err)
	}

	return &m, nil
}

func (r *MemberRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM members WHERE id = $1`

	result, err := r.DB.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete member: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("member not found")
	}

	return nil
}
