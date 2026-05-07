package models

import "time"

type LoanReport struct {
	ID         int        `json:"id"`
	MemberID   int        `json:"member_id"`
	MemberName string     `json:"member_name"`
	CopyID     int        `json:"copy_id"`
	BookTitle  string     `json:"book_title"`
	BookAuthor string     `json:"book_author"`
	BranchID   int        `json:"branch_id"`
	BranchName string     `json:"branch_name"`
	BorrowedAt time.Time  `json:"borrowed_at"`
	DueAt      time.Time  `json:"due_at"`
	ReturnedAt *time.Time `json:"returned_at"`
	Status     string     `json:"status"`
}
