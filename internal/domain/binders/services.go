package binders

import (
	"context"

	"github.com/google/uuid"
	"github.com/nextlevellabs/go-wise/wise"
	"github.com/seosoojin/dalkom/internal/domain/pagination"
	"github.com/seosoojin/dalkom/pkg/models"
)

type Service interface {
	Create(ctx context.Context, binder *models.Binder) error

	GetByUserID(ctx context.Context, userID string, pagination pagination.Page) ([]models.Binder, error)
	GetByID(ctx context.Context, id string) (models.Binder, error)

	Update(ctx context.Context, binder *models.Binder) error

	Delete(ctx context.Context, id string) (models.Binder, error)
}

type service struct {
	repo Repository
}

var _ Service = &service{}

func NewService(repo Repository) *service {
	return &service{
		repo: repo,
	}
}

func (s *service) Create(ctx context.Context, binder *models.Binder) error {
	binder.ID = uuid.NewString()
	return s.repo.Upsert(ctx, binder.ID, *binder)
}

func (s *service) GetByUserID(ctx context.Context, userID string, pagination pagination.Page) ([]models.Binder, error) {
	filter := map[string][]interface{}{
		"user_id": {userID},
	}

	return s.repo.Search(ctx, filter, wise.WithPage(pagination.Offset), wise.WithPageSize(pagination.Limit))
}

func (s *service) GetByID(ctx context.Context, id string) (models.Binder, error) {
	return s.repo.FindOne(ctx, id)
}

func (s *service) Update(ctx context.Context, binder *models.Binder) error {
	_, err := s.GetByID(ctx, binder.ID)
	if err != nil {
		return err
	}
	return s.repo.Upsert(ctx, binder.ID, *binder)
}

func (s *service) Delete(ctx context.Context, id string) (models.Binder, error) {
	return s.repo.Delete(ctx, id)
}
