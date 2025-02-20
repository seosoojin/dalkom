package collections

import (
	"github.com/nextlevellabs/go-wise/wise"
	"github.com/seosoojin/dalkom/pkg/models"
)

type Repository wise.MongoRepository[models.Collection]
