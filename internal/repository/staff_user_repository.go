package repository

import (
	"context"
	"database/sql"
	"fmt"

	"test1/internal/models"
)

type StaffUserRepository struct {
	DB *sql.DB
}

func NewStaffUserRepository(db *sql.DB) *StaffUserRepository {
	return &StaffUserRepository{DB: db}
}

func (r *StaffUserRepository) Create(ctx context.Context, accountID int, firstName, lastName, position string, branchID *int) (*models.StaffUser, error) {
	query := `
		INSERT INTO staff_users (account_id, first_name, last_name, position, branch_id)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, account_id, first_name, last_name, position, branch_id, created_at
	`

	var s models.StaffUser
	err := r.DB.QueryRowContext(ctx, query, accountID, firstName, lastName, position, branchID).
		Scan(&s.ID, &s.AccountID, &s.FirstName, &s.LastName, &s.Position, &s.BranchID, &s.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("create staff user: %w", err)
	}

	return &s, nil
}

func (r *StaffUserRepository) GetAuthUserByEmail(ctx context.Context, email string) (*models.AuthUserDTO, error) {
	query := `
		SELECT a.id, a.email, a.role, s.first_name, s.last_name, s.position, s.branch_id
		FROM accounts a
		JOIN staff_users s ON s.account_id = a.id
		WHERE a.email = $1
	`

	var u models.AuthUserDTO
	err := r.DB.QueryRowContext(ctx, query, email).
		Scan(&u.AccountID, &u.Email, &u.Role, &u.FirstName, &u.LastName, &u.Position, &u.BranchID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("staff user not found")
		}
		return nil, fmt.Errorf("get auth user by email: %w", err)
	}

	return &u, nil
}

func (r *StaffUserRepository) GetAuthUserByAccountID(ctx context.Context, accountID int) (*models.AuthUserDTO, error) {
	query := `
		SELECT a.id, a.email, a.role, s.first_name, s.last_name, s.position, s.branch_id
		FROM accounts a
		JOIN staff_users s ON s.account_id = a.id
		WHERE a.id = $1
	`

	var u models.AuthUserDTO
	err := r.DB.QueryRowContext(ctx, query, accountID).
		Scan(&u.AccountID, &u.Email, &u.Role, &u.FirstName, &u.LastName, &u.Position, &u.BranchID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("staff user not found")
		}
		return nil, fmt.Errorf("get auth user by account id: %w", err)
	}

	return &u, nil
}

func (r *StaffUserRepository) GetByBranch(ctx context.Context, branchID int) ([]models.StaffUser, error) {
	query := `
		SELECT id, account_id, first_name, last_name, position, branch_id, created_at
		FROM staff_users
		WHERE branch_id = $1
		ORDER BY id ASC
	`

	rows, err := r.DB.QueryContext(ctx, query, branchID)
	if err != nil {
		return nil, fmt.Errorf("query staff by branch: %w", err)
	}
	defer rows.Close()

	var staff []models.StaffUser

	for rows.Next() {
		var s models.StaffUser
		if err := rows.Scan(
			&s.ID,
			&s.AccountID,
			&s.FirstName,
			&s.LastName,
			&s.Position,
			&s.BranchID,
			&s.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan staff by branch: %w", err)
		}
		staff = append(staff, s)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate staff by branch: %w", err)
	}

	return staff, nil
}
