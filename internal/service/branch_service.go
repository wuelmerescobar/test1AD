package service

import (
	"context"
	"fmt"
	"strings"

	"test1/internal/models"
	"test1/internal/repository"
)

type BranchService struct {
	Repo *repository.BranchRepository
}

func NewBranchService(repo *repository.BranchRepository) *BranchService {
	return &BranchService{Repo: repo}
}

func (s *BranchService) GetAll(ctx context.Context) ([]models.Branch, error) {
	return s.Repo.GetAll(ctx)
}

func (s *BranchService) Create(ctx context.Context, req models.CreateBranchRequest) (*models.Branch, error) {
	req.Name = strings.TrimSpace(req.Name)
	req.Code = strings.TrimSpace(req.Code)
	req.Address = strings.TrimSpace(req.Address)

	if req.Name == "" {
		return nil, fmt.Errorf("name is required")
	}

	if req.Code == "" {
		return nil, fmt.Errorf("code is required")
	}

	return s.Repo.Create(ctx, req)
}
