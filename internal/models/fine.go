package models

import "time"

type FineReport struct {
	ID         int       `json:"id"`
	LoanID     int       `json:"loan_id"`
	MemberID   int       `json:"member_id"`
	MemberName string    `json:"member_name"`
	CopyID     int       `json:"copy_id"`
	BookTitle  string    `json:"book_title"`
	BranchName string    `json:"branch_name"`
	Amount     string    `json:"amount"`
	Reason     string    `json:"reason"`
	Paid       bool      `json:"paid"`
	CreatedAt  time.Time `json:"created_at"`
}
