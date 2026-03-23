package models

import "time"

type Branch struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Code      string    `json:"code"`
	Address   string    `json:"address"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateBranchRequest struct {
	Name    string `json:"name"`
	Code    string `json:"code"`
	Address string `json:"address"`
}
