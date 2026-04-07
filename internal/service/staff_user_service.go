package service

import (
	"context"
	"fmt"

	"test1/internal/models"
	"test1/internal/repository"
)

type StaffUserService struct {
	Repo *repository.StaffUserRepository
}

func NewStaffUserService(repo *repository.StaffUserRepository) *StaffUserService {
	return &StaffUserService{Repo: repo}
}

func (s *StaffUserService) GetByBranch(ctx context.Context, branchID int) ([]models.StaffUser, error) {
	if branchID <= 0 {
		return nil, fmt.Errorf("invalid branch id")
	}
	return s.Repo.GetByBranch(ctx, branchID)
}
