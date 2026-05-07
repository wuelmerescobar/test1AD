package repository

import (
	"context"
	"database/sql"
	"fmt"

	"test1/internal/models"
)

type LoanRepository struct {
	DB *sql.DB
}

func NewLoanRepository(db *sql.DB) *LoanRepository {
	return &LoanRepository{DB: db}
}

func (r *LoanRepository) GetAllReports(ctx context.Context) ([]models.LoanReport, error) {
	query := `
		SELECT
			l.id,
			m.id,
			m.first_name || ' ' || m.last_name AS member_name,
			bc.id,
			b.title,
			b.author,
			br.id,
			br.name,
			l.borrowed_at,
			l.due_at,
			l.returned_at,
			l.status
		FROM loans l
		JOIN members m ON l.member_id = m.id
		JOIN book_copies bc ON l.book_copy_id = bc.id
		JOIN books b ON bc.book_id = b.id
		JOIN branches br ON bc.branch_id = br.id
		ORDER BY l.borrowed_at DESC, l.id DESC
	`

	rows, err := r.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query loan reports: %w", err)
	}
	defer rows.Close()

	var loans []models.LoanReport
	for rows.Next() {
		var loan models.LoanReport
		if err := rows.Scan(
			&loan.ID,
			&loan.MemberID,
			&loan.MemberName,
			&loan.CopyID,
			&loan.BookTitle,
			&loan.BookAuthor,
			&loan.BranchID,
			&loan.BranchName,
			&loan.BorrowedAt,
			&loan.DueAt,
			&loan.ReturnedAt,
			&loan.Status,
		); err != nil {
			return nil, fmt.Errorf("scan loan report: %w", err)
		}
		loans = append(loans, loan)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate loan reports: %w", err)
	}

	return loans, nil
}
