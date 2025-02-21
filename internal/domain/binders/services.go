package binders

import (
	"context"

	"github.com/google/uuid"
	"github.com/nextlevellabs/go-wise/wise"
	"github.com/seosoojin/dalkom/internal/domain/cards"
	"github.com/seosoojin/dalkom/internal/domain/pagination"
	"github.com/seosoojin/dalkom/pkg/models"
)

type Service interface {
	Create(ctx context.Context, binder *models.Binder) error

	GetByUserID(ctx context.Context, userID string, filter map[string][]any, pagination pagination.Page) ([]models.Binder, error)
	GetByID(ctx context.Context, id string) (models.Binder, error)
	GetBinderCards(ctx context.Context, id string) ([]models.Card, error)

	Update(ctx context.Context, binder *models.Binder) error
	AddCard(ctx context.Context, binderID, cardID string) error
	RemoveCard(ctx context.Context, binderID, cardID string) error

	Delete(ctx context.Context, id string) (models.Binder, error)
}

type service struct {
	repo     Repository
	cardRepo cards.Repository
}

var _ Service = &service{}

func NewService(repo Repository, cardRepo cards.Repository) *service {
	return &service{
		repo:     repo,
		cardRepo: cardRepo,
	}
}

func (s *service) Create(ctx context.Context, binder *models.Binder) error {
	binder.ID = uuid.NewString()
	return s.repo.Upsert(ctx, binder.ID, *binder)
}

func (s *service) GetByUserID(ctx context.Context, userID string, filter map[string][]any, pagination pagination.Page) ([]models.Binder, error) {
	queryFilter := map[string][]interface{}{
		"user_id": {userID},
	}

	for k, v := range filter {
		queryFilter[k] = v
	}

	return s.repo.Search(ctx, queryFilter, wise.WithPage(pagination.Offset), wise.WithPageSize(pagination.Limit))
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

func (s *service) AddCard(ctx context.Context, binderID, cardID string) error {
	binder, err := s.GetByID(ctx, binderID)
	if err != nil {
		return err
	}

	if binder.CardIDs == nil {
		binder.CardIDs = &[]string{}
	}

	*binder.CardIDs = append(*binder.CardIDs, cardID)

	return s.repo.Upsert(ctx, binder.ID, binder)
}

func (s *service) RemoveCard(ctx context.Context, binderID, cardID string) error {
	binder, err := s.GetByID(ctx, binderID)
	if err != nil {
		return err
	}

	if binder.CardIDs == nil {
		binder.CardIDs = &[]string{}
	}

	for i, id := range *binder.CardIDs {
		if id == cardID {
			*binder.CardIDs = append((*binder.CardIDs)[:i], (*binder.CardIDs)[i+1:]...)
			break
		}
	}

	return s.repo.Upsert(ctx, binder.ID, binder)
}

func (s *service) GetBinderCards(ctx context.Context, id string) ([]models.Card, error) {
	binder, err := s.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	cards, err := s.cardRepo.Find(ctx, *binder.CardIDs)
	if err != nil {
		return nil, err
	}

	cardsResultMap := make(map[string]models.Card, len(cards))
	for _, card := range cards {
		cardsResultMap[card.ID] = card
	}

	result := make([]models.Card, 0, len(*binder.CardIDs))
	for _, id := range *binder.CardIDs {
		result = append(result, cardsResultMap[id])
	}

	return result, nil
}

func (s *service) Delete(ctx context.Context, id string) (models.Binder, error) {
	return s.repo.Delete(ctx, id)
}
