package users

import (
	"context"
	"net/mail"
	"strings"

	"github.com/google/uuid"
	"github.com/seosoojin/dalkom/internal/domain/auth"
	"github.com/seosoojin/dalkom/pkg/models"
)

type Service interface {
	Login(ctx context.Context, user, password string) (string, error)
	RefreshToken(ctx context.Context) (string, error)
	GetByID(ctx context.Context, id string) (models.User, error)
	GetByEmail(ctx context.Context, email string) (models.User, error)
	Create(ctx context.Context, user *models.User) error
	Update(ctx context.Context, user *models.User) error
}

type service struct {
	repo   Repository
	auth   auth.Services
	jwtSvc auth.JWTService
}

var _ Service = &service{}

func NewService(repo Repository, jwtSvc auth.JWTService) *service {
	return &service{
		repo:   repo,
		auth:   auth.NewService(),
		jwtSvc: jwtSvc,
	}
}

func (s *service) Login(ctx context.Context, user, password string) (string, error) {
	user = strings.ToLower(user)

	filter := map[string][]any{
		"email": {user},
	}

	u, err := s.repo.Search(ctx, filter)
	if err != nil {
		return "", err
	}

	if len(u) == 0 {
		return "", ErrInvalidCredentials
	}

	if err := s.auth.Verify(u[0].Password, password); err != nil {
		return "", ErrInvalidCredentials
	}

	u[0].Password = ""

	token, err := s.jwtSvc.GenerateToken(u[0])
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *service) RefreshToken(ctx context.Context) (string, error) {
	user := auth.UserFromContext(ctx)

	return s.jwtSvc.GenerateToken(user)
}

func (s *service) GetByID(ctx context.Context, id string) (models.User, error) {
	return s.repo.FindOne(ctx, id)
}

func (s *service) GetByEmail(ctx context.Context, email string) (models.User, error) {
	filter := map[string][]any{
		"email": {email},
	}

	u, err := s.repo.Search(ctx, filter)
	if err != nil {
		return models.User{}, err
	}

	if len(u) == 0 {
		return models.User{}, ErrNotFound
	}

	return u[0], nil
}

func (s *service) Create(ctx context.Context, user *models.User) error {
	user.Username = strings.TrimSpace(user.Username)
	user.Email = strings.TrimSpace(user.Email)

	if err := s.validateUser(user); err != nil {
		return err
	}

	user.Username = strings.ToLower(user.Username)
	user.Email = strings.ToLower(user.Email)

	_, err := s.GetByEmail(ctx, user.Email)
	if err == nil {
		return ErrUserExists
	}

	user.ID = uuid.NewString()
	hashed, err := s.auth.HashPassword(user.Password)
	if err != nil {
		return err
	}

	user.Password = hashed

	return s.repo.Upsert(ctx, user.ID, *user)
}

func (s *service) validateUser(user *models.User) error {

	_, err := mail.ParseAddress(user.Email)
	if err != nil {
		return ErrInvalidEmail
	}

	if len(user.Username) < 4 {
		return ErrInvalidUsername
	}

	if len(user.Password) < 8 {
		return ErrInvalidPassword
	}
	return nil
}

func (s *service) Update(ctx context.Context, user *models.User) error {
	user.Username = strings.TrimSpace(user.Username)
	user.Email = strings.TrimSpace(user.Email)

	if err := s.validateUser(user); err != nil {
		return err
	}

	user.Username = strings.ToLower(user.Username)
	user.Email = strings.ToLower(user.Email)

	_, err := s.GetByID(ctx, user.ID)
	if err != nil {
		return err
	}

	if user.Password != "" {
		hashed, err := s.auth.HashPassword(user.Password)
		if err != nil {
			return err
		}
		user.Password = hashed
	}

	return s.repo.Upsert(ctx, user.ID, *user)
}
