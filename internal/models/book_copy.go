package models

import "time"

type BookCopy struct {
	ID        int       `json:"id"`
	BookID    int       `json:"book_id"`
	BranchID  int       `json:"branch_id"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateBookCopyRequest struct {
	BookID   int    `json:"book_id"`
	BranchID int    `json:"branch_id"`
	Status   string `json:"status"`
}
