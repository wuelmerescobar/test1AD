package models

import "time"

type StaffUser struct {
	ID        int       `json:"id"`
	AccountID int       `json:"account_id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Position  string    `json:"position"`
	BranchID  *int      `json:"branch_id"`
	CreatedAt time.Time `json:"created_at"`
}
