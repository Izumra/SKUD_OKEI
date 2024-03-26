package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/Izumra/SKUD_OKEI/domain/entity"
	"github.com/Izumra/SKUD_OKEI/internal/storage"
)

func (s *Storage) UserByID(ctx context.Context, id int64) (*entity.User, error) {
	op := "sqlite/UserStorage.UserByID"

	query := "select * from users where id=?"
	state, err := s.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	results, err := state.QueryContext(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, storage.ErrUserNotFound
		}
		return nil, err
	}

	var user entity.User
	if !results.Next() {
		return nil, storage.ErrUserNotFound
	}

	err = results.Scan(&user.Id, &user.Username, &user.Password, &user.Role)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &user, nil
}

func (s *Storage) UserByUsername(ctx context.Context, username string) (*entity.User, error) {
	op := "sqlite/UserStorage.UserByUsername"

	query := "select * from users where username=?"
	state, err := s.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	results, err := state.QueryContext(ctx, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, storage.ErrUserNotFound
		}
		return nil, err
	}

	var user entity.User
	if !results.Next() {
		return nil, storage.ErrUserNotFound
	}

	err = results.Scan(&user.Id, &user.Username, &user.Password, &user.Role)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &user, nil
}

func (s *Storage) AddUser(ctx context.Context, data entity.User) (int64, error) {
	op := "storage/sqlite/UserStorage.AddUser"

	query := "insert into users(username,pass,role)values(?,?,?)"
	state, err := s.db.PrepareContext(ctx, query)
	if err != nil {
		return -1, fmt.Errorf("%s: %w", op, err)
	}

	result, err := state.ExecContext(ctx, data.Username, data.Password, data.Role)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return -1, storage.ErrUserExist
		}
		return -1, fmt.Errorf("%s: %w", op, err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return -1, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}
func (s *Storage) DeleteUserById(ctx context.Context, id int64) error {
	op := "storage/sqlite/UserStorage.DeleteUserById"

	query := "delete from users where id=?"
	state, err := s.db.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = state.ExecContext(ctx, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
