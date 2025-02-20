package groups

import (
	"context"

	"github.com/google/uuid"
	"github.com/seosoojin/dalkom/pkg/models"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type Service interface {
	GetGroups(ctx context.Context) ([]models.Group, error)

	Create(ctx context.Context, group *models.Group) error

	Update(ctx context.Context, group *models.Group) error

	Delete(ctx context.Context, id string) (models.Group, error)
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

func (s *service) GetGroups(ctx context.Context) ([]models.Group, error) {
	return s.repo.FindAll(ctx)
}

func (s *service) Create(ctx context.Context, group *models.Group) error {
	group.ID = uuid.NewString()
	group.Name = s.caser.String(group.Name)
	return s.repo.Upsert(ctx, group.ID, *group)
}

func (s *service) Update(ctx context.Context, group *models.Group) error {
	if group.Name != "" {
		group.Name = s.caser.String(group.Name)
	}
	return s.repo.Upsert(ctx, group.ID, *group)
}

func (s *service) Delete(ctx context.Context, id string) (models.Group, error) {
	return s.repo.Delete(ctx, id)
}
