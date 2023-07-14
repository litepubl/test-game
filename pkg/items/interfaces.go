// Package items implements application business logic. Each logic group in own file.
package items

import (
	"context"

	"github.com/litepubl/test-game/pkg/entity"
)

//go:generate mockgen -source=interfaces.go -destination=./mocks_test.go -package=items_test
type (
	// CRUDItems интерфейс  CRUD для Itms
	CRUDItems interface {
		List(ctx context.Context) ([]byte, error)
		Create(ctx context.Context, campaignId int, name string) (entity.Item, error)
		Update(ctx context.Context, u entity.UpdateData) (entity.Item, error)
		Delete(ctx context.Context, id int) error
	}

	// ItemsRepo  репозиторий для items
	ItemsRepo interface {
		Create(ctx context.Context, item *entity.Item) error
		Update(ctx context.Context, u entity.UpdateData) (entity.Item, error)
		Delete(ctx context.Context, id int) (entity.Item, error)
		List(ctx context.Context) ([]entity.Item, error)
		MaxPriority(ctx context.Context) (int, error)
	}
)
