package service

import (
	"context"

	"test1/internal/models"
	"test1/internal/repository"
)

type FineService struct {
	Repo *repository.FineRepository
}

func NewFineService(repo *repository.FineRepository) *FineService {
	return &FineService{Repo: repo}
}

func (s *FineService) GetAllReports(ctx context.Context) ([]models.FineReport, error) {
	return s.Repo.GetAllReports(ctx)
}
