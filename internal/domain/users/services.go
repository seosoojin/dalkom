package users

import (
	"context"

	"github.com/google/uuid"
	"github.com/seosoojin/dalkom/internal/domain/auth"
	"github.com/seosoojin/dalkom/pkg/models"
)

type Service interface {
	Login(ctx context.Context, user, password string) (string, error)
	GetByID(ctx context.Context, id string) (models.User, error)
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
	filter := map[string][]any{
		"$or": {
			map[string]any{"email": user},
			map[string]any{"username": user},
		},
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

	token, err := s.jwtSvc.GenerateToken(u[0])
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *service) GetByID(ctx context.Context, id string) (models.User, error) {
	return s.repo.FindOne(ctx, id)
}

func (s *service) Create(ctx context.Context, user *models.User) error {
	user.ID = uuid.NewString()
	return s.repo.Upsert(ctx, user.ID, *user)
}

func (s *service) Update(ctx context.Context, user *models.User) error {
	_, err := s.GetByID(ctx, user.ID)
	if err != nil {
		return err
	}
	return s.repo.Upsert(ctx, user.ID, *user)
}
