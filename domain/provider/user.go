package provider

import (
	"context"

	"github.com/Izumra/SKUD_OKEI/domain/entity"
)

type User interface {
	UserByID(ctx context.Context, id int64) (*entity.User, error)
	UserByUsername(ctx context.Context, username string) (*entity.User, error)
}
