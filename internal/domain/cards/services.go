package cards

import (
	"context"

	"github.com/google/uuid"
	"github.com/seosoojin/dalkom/internal/domain/pagination"
	"github.com/seosoojin/dalkom/pkg/models"
)

type Service interface {
	GetCard(ctx context.Context, id string) (models.Card, error)
	GetCards(ctx context.Context, pagination pagination.Page) ([]models.Card, error)

	CreateCard(ctx context.Context, card *models.Card) error

	UpdateCard(ctx context.Context, card *models.Card) error

	DeleteCard(ctx context.Context, id string) (models.Card, error)
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

func (s *service) GetCard(ctx context.Context, id string) (models.Card, error) {
	return s.repo.FindOne(ctx, id)
}

func (s *service) GetCards(ctx context.Context, pagination pagination.Page) ([]models.Card, error) {
	return s.repo.Search(ctx, nil)
}

func (s *service) CreateCard(ctx context.Context, card *models.Card) error {
	card.ID = uuid.NewString()
	return s.repo.Upsert(ctx, card.ID, *card)
}

func (s *service) UpdateCard(ctx context.Context, card *models.Card) error {
	return s.repo.Upsert(ctx, card.ID, *card)
}

func (s *service) DeleteCard(ctx context.Context, id string) (models.Card, error) {
	return s.repo.Delete(ctx, id)
}
