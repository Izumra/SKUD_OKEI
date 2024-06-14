package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/Izumra/SKUD_OKEI/domain/dto/resp"
	"github.com/Izumra/SKUD_OKEI/domain/entity"
	"github.com/Izumra/SKUD_OKEI/domain/provider"
	"github.com/Izumra/SKUD_OKEI/domain/repository"
	valueobject "github.com/Izumra/SKUD_OKEI/domain/value-object"
	"github.com/Izumra/SKUD_OKEI/internal/storage"
)

var (
	ErrUserAlreadyRegistered = errors.New("пользователь с такими данными уже зарегестрирован в системе")
)

type SessionStorage interface {
	Create(ctx context.Context, data *entity.User) (sessionId string, err error)
	GetByID(ctx context.Context, sessionId string) (*entity.User, error)
	DeleteByID(ctx context.Context, sessionId string) error
	UpdateByID(ctx context.Context, sessionId string, updatedData *entity.User) error
}

type Service struct {
	logger      *slog.Logger
	sessStorage SessionStorage
	usrRep      repository.User
	usrPrvdr    provider.User
}

func NewService(
	logger *slog.Logger,
	sessStorage SessionStorage,
	usrRep repository.User,
	usrPrvdr provider.User,
) *Service {
	return &Service{
		logger,
		sessStorage,
		usrRep,
		usrPrvdr,
	}
}

func (s *Service) Login(ctx context.Context, username, password string) (*resp.SuccessAuth, error) {
	op := "internal/services/auth.Service.Login"
	logger := s.logger.With(slog.String("op", op))

	user, err := s.usrPrvdr.UserByUsername(ctx, username)
	if err != nil {
		logger.Error("Occured the error while finding the user", err)
		return nil, err
	}
	if user.Password != password {
		return nil, fmt.Errorf("Пароли не совпадают")
	}

	sessionId, err := s.sessStorage.Create(ctx, user)
	if err != nil {
		logger.Error("Occured the error while creating the session", err)
		return nil, err
	}

	return &resp.SuccessAuth{
		Username:  user.Username,
		SessionId: sessionId,
	}, nil
}

func (s *Service) Registrate(ctx context.Context, username, password string) (*resp.SuccessAuth, error) {
	op := "internal/services/auth.Service.Registrate"
	logger := s.logger.With(slog.String("op", op))

	user := entity.User{
		Username: username,
		Password: password,
		Role:     valueobject.StudentRole,
	}

	userId, err := s.usrRep.AddUser(ctx, user)
	if err != nil {
		if errors.Is(err, storage.ErrUserExist) {
			return nil, ErrUserAlreadyRegistered
		}
		logger.Error("Occured the error while finding the user", err)
		return nil, err
	}
	user.Id = userId

	sessionId, err := s.sessStorage.Create(ctx, &user)
	if err != nil {
		logger.Error("Occured the error while creating the session", err)
		return nil, err
	}

	return &resp.SuccessAuth{
		Username:  user.Username,
		SessionId: sessionId,
	}, nil
}
