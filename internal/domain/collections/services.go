package collections

import (
	"context"

	"github.com/google/uuid"
	"github.com/seosoojin/dalkom/pkg/models"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type Service interface {
	GetCollections(context.Context) ([]models.Collection, error)

	Create(ctx context.Context, collection *models.Collection) error

	Update(ctx context.Context, collection *models.Collection) error

	Delete(ctx context.Context, id string) (models.Collection, error)
}

type service struct {
	repo  Repository
	caser cases.Caser
}

var _ Service = &service{}

func NewService(repo Repository) *service {
	return &service{
		repo:  repo,
		caser: cases.Title(language.English),
	}
}

func (s *service) GetCollections(ctx context.Context) ([]models.Collection, error) {
	return s.repo.FindAll(ctx)
}

func (s *service) Create(ctx context.Context, collection *models.Collection) error {
	collection.ID = uuid.NewString()
	collection.Name = s.caser.String(collection.Name)
	return s.repo.Upsert(ctx, collection.ID, *collection)
}

func (s *service) Update(ctx context.Context, collection *models.Collection) error {
	if collection.Name != "" {
		collection.Name = s.caser.String(collection.Name)
	}
	return s.repo.Upsert(ctx, collection.ID, *collection)
}

func (s *service) Delete(ctx context.Context, id string) (models.Collection, error) {
	return s.repo.Delete(ctx, id)
}
