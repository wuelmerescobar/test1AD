package service

import (
	"context"

	"test1/internal/models"
	"test1/internal/repository"
)

type LoanService struct {
	Repo *repository.LoanRepository
}

func NewLoanService(repo *repository.LoanRepository) *LoanService {
	return &LoanService{Repo: repo}
}

func (s *LoanService) GetAllReports(ctx context.Context) ([]models.LoanReport, error) {
	return s.Repo.GetAllReports(ctx)
}
