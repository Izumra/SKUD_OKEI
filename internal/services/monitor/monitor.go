package monitor

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/Izumra/SKUD_OKEI/internal/services/auth"
)

var (
	ErrSessionTokenInvalid = errors.New("сессия пользователя не действительна")
	ErrAccessDenied        = errors.New("вам отказано в доступе")
)

type Service struct {
	logger         *slog.Logger
	sessStore      auth.SessionStorage
	integrServAddr string
}

func NewService(
	logger *slog.Logger,
	sessStore auth.SessionStorage,
	integrServAddr string,
) *Service {
	return &Service{
		logger,
		sessStore,
		fmt.Sprintf("%s/soap/IOrionPro", integrServAddr),
	}
}
