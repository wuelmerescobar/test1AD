package service

import (
	"context"
	"fmt"
	"strings"

	"test1/internal/models"
	"test1/internal/repository"
)

type BookCopyService struct {
	Repo *repository.BookCopyRepository
}

func NewBookCopyService(repo *repository.BookCopyRepository) *BookCopyService {
	return &BookCopyService{Repo: repo}
}

func (s *BookCopyService) Create(ctx context.Context, req models.CreateBookCopyRequest) (*models.BookCopy, error) {
	if req.BookID <= 0 {
		return nil, fmt.Errorf("book_id is required")
	}

	if req.BranchID <= 0 {
		return nil, fmt.Errorf("branch_id is required")
	}

	req.Status = strings.TrimSpace(req.Status)
	if req.Status == "" {
		req.Status = "available"
	}

	return s.Repo.Create(ctx, req)
}
