package repository

import (
	"context"

	"github.com/Izumra/SKUD_OKEI/domain/entity"
)

type User interface {
	AddUser(ctx context.Context, data entity.User) (int64, error)
	DeleteUserById(ctx context.Context, id int64) error
}
