package repository

import (
	"context"
	"database/sql"
	"fmt"

	"test1/internal/models"
)

type FineRepository struct {
	DB *sql.DB
}

func NewFineRepository(db *sql.DB) *FineRepository {
	return &FineRepository{DB: db}
}

func (r *FineRepository) GetAllReports(ctx context.Context) ([]models.FineReport, error) {
	query := `
		SELECT
			f.id,
			l.id,
			m.id,
			m.first_name || ' ' || m.last_name AS member_name,
			bc.id,
			b.title,
			br.name,
			f.amount::text,
			f.reason,
			f.paid,
			f.created_at
		FROM fines f
		JOIN loans l ON f.loan_id = l.id
		JOIN members m ON f.member_id = m.id
		JOIN book_copies bc ON l.book_copy_id = bc.id
		JOIN books b ON bc.book_id = b.id
		JOIN branches br ON bc.branch_id = br.id
		ORDER BY f.created_at DESC, f.id DESC
	`

	rows, err := r.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query fine reports: %w", err)
	}
	defer rows.Close()

	var fines []models.FineReport
	for rows.Next() {
		var fine models.FineReport
		if err := rows.Scan(
			&fine.ID,
			&fine.LoanID,
			&fine.MemberID,
			&fine.MemberName,
			&fine.CopyID,
			&fine.BookTitle,
			&fine.BranchName,
			&fine.Amount,
			&fine.Reason,
			&fine.Paid,
			&fine.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan fine report: %w", err)
		}
		fines = append(fines, fine)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate fine reports: %w", err)
	}

	return fines, nil
}
