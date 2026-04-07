package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"test1/internal/mailer"
	"test1/internal/models"
	"test1/internal/repository"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	AccountRepo *repository.AccountRepository
	StaffRepo   *repository.StaffUserRepository
	Mailer      *mailer.Mailer
	JWTSecret   string
}

func NewAuthService(accountRepo *repository.AccountRepository, staffRepo *repository.StaffUserRepository, mailer *mailer.Mailer, jwtSecret string) *AuthService {
	return &AuthService{
		AccountRepo: accountRepo,
		StaffRepo:   staffRepo,
		Mailer:      mailer,
		JWTSecret:   jwtSecret,
	}
}

func (s *AuthService) RegisterStaff(ctx context.Context, req models.RegisterStaffRequest) (*models.AuthUserDTO, error) {
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	req.Password = strings.TrimSpace(req.Password)
	req.Role = strings.TrimSpace(strings.ToLower(req.Role))
	req.FirstName = strings.TrimSpace(req.FirstName)
	req.LastName = strings.TrimSpace(req.LastName)
	req.Position = strings.TrimSpace(req.Position)

	if req.Email == "" {
		return nil, fmt.Errorf("email is required")
	}
	if req.Password == "" || len(req.Password) < 6 {
		return nil, fmt.Errorf("password must be at least 6 characters")
	}
	if req.FirstName == "" {
		return nil, fmt.Errorf("first_name is required")
	}
	if req.LastName == "" {
		return nil, fmt.Errorf("last_name is required")
	}
	if req.Role == "" {
		req.Role = "viewer"
	}
	if req.Role != "admin" && req.Role != "librarian" && req.Role != "viewer" {
		return nil, fmt.Errorf("invalid role")
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	account, err := s.AccountRepo.Create(ctx, req.Email, string(passwordHash), req.Role)
	if err != nil {
		return nil, err
	}

	_, err = s.StaffRepo.Create(ctx, account.ID, req.FirstName, req.LastName, req.Position, req.BranchID)
	if err != nil {
		return nil, err
	}

	_ = s.Mailer.SendWelcomeStaffEmail(req.Email, req.FirstName+" "+req.LastName)

	return s.StaffRepo.GetAuthUserByAccountID(ctx, account.ID)
}

func (s *AuthService) Login(ctx context.Context, req models.LoginRequest) (*models.LoginResponse, error) {
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	req.Password = strings.TrimSpace(req.Password)

	if req.Email == "" || req.Password == "" {
		return nil, fmt.Errorf("email and password are required")
	}

	account, err := s.AccountRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	if !account.IsActive {
		return nil, fmt.Errorf("account is inactive")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(account.PasswordHash), []byte(req.Password)); err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	user, err := s.StaffRepo.GetAuthUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	token, err := s.GenerateJWT(account.ID, account.Email, account.Role)
	if err != nil {
		return nil, err
	}

	return &models.LoginResponse{
		Token: token,
		User:  *user,
	}, nil
}

func (s *AuthService) GenerateJWT(accountID int, email, role string) (string, error) {
	claims := jwt.MapClaims{
		"sub":   accountID,
		"email": email,
		"role":  role,
		"exp":   time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.JWTSecret))
}

func (s *AuthService) GetMe(ctx context.Context, accountID int) (*models.AuthUserDTO, error) {
	return s.StaffRepo.GetAuthUserByAccountID(ctx, accountID)
}
