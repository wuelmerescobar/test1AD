package service

import (
	"context"
	"fmt"
	"strings"

	"test1/internal/models"
	"test1/internal/repository"
)

type BookService struct {
	Repo *repository.BookRepository
}

func NewBookService(repo *repository.BookRepository) *BookService {
	return &BookService{Repo: repo}
}

func (s *BookService) GetAll(ctx context.Context) ([]models.Book, error) {
	return s.Repo.GetAll(ctx)
}

func (s *BookService) Create(ctx context.Context, req models.CreateBookRequest) (*models.Book, error) {
	req.Title = strings.TrimSpace(req.Title)
	req.Author = strings.TrimSpace(req.Author)
	req.ISBN = strings.TrimSpace(req.ISBN)
	req.Genre = strings.TrimSpace(req.Genre)

	if req.Title == "" {
		return nil, fmt.Errorf("title is required")
	}

	if req.Author == "" {
		return nil, fmt.Errorf("author is required")
	}

	return s.Repo.Create(ctx, req)
}

func (s *BookService) GetAvailableByBranch(ctx context.Context, branchID int) ([]models.Book, error) {
	if branchID <= 0 {
		return nil, fmt.Errorf("invalid branch id")
	}

	return s.Repo.GetAvailableByBranch(ctx, branchID)
}
func (s *BookService) Update(ctx context.Context, id int, req models.CreateBookRequest) (*models.Book, error) {
	req.Title = strings.TrimSpace(req.Title)
	req.Author = strings.TrimSpace(req.Author)
	req.ISBN = strings.TrimSpace(req.ISBN)
	req.Genre = strings.TrimSpace(req.Genre)

	if id <= 0 {
		return nil, fmt.Errorf("invalid book id")
	}
	if req.Title == "" {
		return nil, fmt.Errorf("title is required")
	}
	if req.Author == "" {
		return nil, fmt.Errorf("author is required")
	}

	return s.Repo.Update(ctx, id, req)
}

func (s *BookService) Delete(ctx context.Context, id int) error {
	if id <= 0 {
		return fmt.Errorf("invalid book id")
	}
	return s.Repo.Delete(ctx, id)
}
