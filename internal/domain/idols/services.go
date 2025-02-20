package idols

import (
	"context"

	"github.com/google/uuid"
	"github.com/seosoojin/dalkom/pkg/models"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type Service interface {
	GetIdols(context.Context) ([]models.Idol, error)

	Create(ctx context.Context, idol *models.Idol) error

	Update(ctx context.Context, idol *models.Idol) error

	Delete(ctx context.Context, id string) (models.Idol, error)
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

func (s *service) GetIdols(ctx context.Context) ([]models.Idol, error) {
	return s.repo.FindAll(ctx)
}

func (s *service) Create(ctx context.Context, idol *models.Idol) error {
	idol.ID = uuid.NewString()
	idol.Name = s.caser.String(idol.Name)
	idol.StageName = s.caser.String(idol.StageName)
	return s.repo.Upsert(ctx, idol.ID, *idol)
}

func (s *service) Update(ctx context.Context, idol *models.Idol) error {
	if idol.Name != "" {
		idol.Name = s.caser.String(idol.Name)
	}

	if idol.StageName != "" {
		idol.StageName = s.caser.String(idol.StageName)
	}

	return s.repo.Upsert(ctx, idol.ID, *idol)
}

func (s *service) Delete(ctx context.Context, id string) (models.Idol, error) {
	return s.repo.Delete(ctx, id)
}
