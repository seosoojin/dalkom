package cards

import (
	"context"

	"github.com/google/uuid"
	"github.com/nextlevellabs/go-wise/wise"
	"github.com/seosoojin/dalkom/internal/domain/collections"
	"github.com/seosoojin/dalkom/internal/domain/groups"
	"github.com/seosoojin/dalkom/internal/domain/idols"
	"github.com/seosoojin/dalkom/internal/domain/pagination"
	"github.com/seosoojin/dalkom/pkg/models"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type Service interface {
	GetCard(ctx context.Context, id string) (models.Card, error)
	GetEnrichedCard(ctx context.Context, id string) (models.EnrichedCard, error)
	GetCards(ctx context.Context, filter map[string][]any, pagination pagination.Page) ([]models.Card, error)

	CreateCard(ctx context.Context, card *models.Card) error

	UpdateCard(ctx context.Context, card *models.Card) error

	DeleteCard(ctx context.Context, id string) (models.Card, error)
}

type service struct {
	repo           Repository
	groupRepo      groups.Repository
	idolRepo       idols.Repository
	collectionRepo collections.Repository
	caser          cases.Caser
}

var _ Service = &service{}

func NewService(repo Repository, groupRepo groups.Repository, idolRepo idols.Repository, collectionRepo collections.Repository) *service {
	return &service{
		repo:           repo,
		groupRepo:      groupRepo,
		idolRepo:       idolRepo,
		collectionRepo: collectionRepo,
		caser:          cases.Title(language.English),
	}
}

func (s *service) GetCard(ctx context.Context, id string) (models.Card, error) {
	return s.repo.FindOne(ctx, id)
}

func (s *service) GetEnrichedCard(ctx context.Context, id string) (models.EnrichedCard, error) {
	card, err := s.GetCard(ctx, id)
	if err != nil {
		return models.EnrichedCard{}, err
	}

	group, err := s.groupRepo.FindOne(ctx, card.GroupID)
	if err != nil {
		return models.EnrichedCard{}, err
	}

	idols, err := s.idolRepo.Find(ctx, card.IdolIDs)
	if err != nil {
		return models.EnrichedCard{}, err
	}

	collection, err := s.collectionRepo.FindOne(ctx, card.CollectionID)
	if err != nil {
		return models.EnrichedCard{}, err
	}

	fmtype := models.ShortTypesMap[card.Type]
	if fmtype == "" {
		fmtype = "R"
	}

	return models.EnrichedCard{
		ID:         card.ID,
		Name:       card.Name,
		ShortName:  card.ShortName,
		Type:       card.Type,
		FmtType:    fmtype,
		ImageUrl:   card.ImageUrl,
		Group:      group,
		Idols:      idols,
		Collection: collection,
	}, nil
}

func (s *service) GetCards(ctx context.Context, filter map[string][]any, pagination pagination.Page) ([]models.Card, error) {
	return s.repo.Search(ctx, filter, wise.WithPage(pagination.Offset), wise.WithPageSize(pagination.Limit), wise.WithSort(map[string]int{"type": 1}))
}

func (s *service) CreateCard(ctx context.Context, card *models.Card) error {
	card.ID = uuid.NewString()
	card.Name = s.caser.String(card.Name)
	return s.repo.Upsert(ctx, card.ID, *card)
}

func (s *service) UpdateCard(ctx context.Context, card *models.Card) error {
	if card.Name != "" {
		card.Name = s.caser.String(card.Name)
	}

	return s.repo.Upsert(ctx, card.ID, *card)
}

func (s *service) DeleteCard(ctx context.Context, id string) (models.Card, error) {
	return s.repo.Delete(ctx, id)
}
