package embedded

import (
	"context"
	"fmt"

	"github.com/Izumra/SKUD_OKEI/domain/entity"
	"github.com/Izumra/SKUD_OKEI/internal/storage/cache"
	"github.com/google/uuid"
)

type SessionStorage struct {
	storage map[string]*entity.User
}

func NewSessStore() *SessionStorage {
	return &SessionStorage{
		storage: make(map[string]*entity.User),
	}
}

func (ss *SessionStorage) Create(ctx context.Context, data *entity.User) (sessionId string, err error) {
	op := "storage/cache/embedded/SessionStorage.Add"

	randID, err := uuid.NewV7()
	if err != nil {
		return "", fmt.Errorf("%s:%w", op, err)
	}

	sessionId = randID.String()
	ss.storage[sessionId] = data

	return sessionId, nil
}

func (ss *SessionStorage) GetByID(ctx context.Context, sessionId string) (*entity.User, error) {
	op := "storage/cache/embedded/SessionStorage.GetByID"

	user, ok := ss.storage[sessionId]
	if !ok {
		return nil, fmt.Errorf("%s: %w", op, cache.ErrSessionNotFound)
	}

	return user, nil
}

func (ss *SessionStorage) DeleteByID(ctx context.Context, sessionId string) error {
	op := "storage/cache/embedded/SessionStorage.DeleteByID"

	_, ok := ss.storage[sessionId]
	if !ok {
		return fmt.Errorf("%s: %w", op, cache.ErrSessionNotFound)
	}

	delete(ss.storage, sessionId)

	return nil
}

func (ss *SessionStorage) UpdateByID(ctx context.Context, sessionId string, updatedData *entity.User) error {
	op := "storage/cache/embedded/SessionStorage.UpdateByID"

	_, ok := ss.storage[sessionId]
	if !ok {
		return fmt.Errorf("%s: %w", op, cache.ErrSessionNotFound)
	}

	ss.storage[sessionId] = updatedData

	return nil
}
