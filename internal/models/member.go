package models

import "time"

type Member struct {
	ID        int       `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	BranchID  *int      `json:"branch_id"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateMemberRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	BranchID  *int   `json:"branch_id"`
}
