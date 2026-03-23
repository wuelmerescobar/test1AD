package service

import (
	"context"
	"fmt"
	"strings"

	"test1/internal/models"
	"test1/internal/repository"
)

type MemberService struct {
	Repo *repository.MemberRepository
}

func NewMemberService(repo *repository.MemberRepository) *MemberService {
	return &MemberService{Repo: repo}
}

func (s *MemberService) GetAll(ctx context.Context) ([]models.Member, error) {
	return s.Repo.GetAll(ctx)
}

func (s *MemberService) GetByBranch(ctx context.Context, branchID int) ([]models.Member, error) {
	if branchID <= 0 {
		return nil, fmt.Errorf("invalid branch id")
	}
	return s.Repo.GetByBranch(ctx, branchID)
}

func (s *MemberService) Create(ctx context.Context, req models.CreateMemberRequest) (*models.Member, error) {
	req.FirstName = strings.TrimSpace(req.FirstName)
	req.LastName = strings.TrimSpace(req.LastName)
	req.Email = strings.TrimSpace(req.Email)
	req.Phone = strings.TrimSpace(req.Phone)

	if req.FirstName == "" {
		return nil, fmt.Errorf("first_name is required")
	}

	if req.LastName == "" {
		return nil, fmt.Errorf("last_name is required")
	}

	return s.Repo.Create(ctx, req)
}

func (s *MemberService) Delete(ctx context.Context, id int) error {
	if id <= 0 {
		return fmt.Errorf("invalid member id")
	}
	return s.Repo.Delete(ctx, id)
}
