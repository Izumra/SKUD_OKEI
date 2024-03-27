package ws

import (
	"context"
	"errors"

	"github.com/Izumra/SKUD_OKEI/domain/dto/integrserv"
	"github.com/Izumra/SKUD_OKEI/domain/dto/reqs"
	"github.com/Izumra/SKUD_OKEI/internal/lib/response"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

var (
	ErrSessionRequired = errors.New(" Для продолжения действия требуется сессия")
	ErrWrongReqBody    = errors.New(" Тело сообщения неверного формата")
)

type MonitorService interface {
	GetEvents(ctx context.Context, sesionId string) (*integrserv.Event, error)
	GetEventsCount(ctx context.Context, sessionId string, eventsFilter *reqs.ReqEventFilter) (int64, error)
}

type MonitorController struct {
	service MonitorService
}

func RegistrMonitorAPI(router fiber.Router, ms MonitorService) {
	mc := MonitorController{
		service: ms,
	}

	router.Use(mc.CheckRegisteredUpgrade())
	router.Get("/events")
	router.Get("/events/count", mc.GetEventsCount())
}

func (mc *MonitorController) CheckRegisteredUpgrade() fiber.Handler {
	return func(c *fiber.Ctx) error {
		sessionId := c.Get("Authorization")
		if websocket.IsWebSocketUpgrade(c) && sessionId != "" {
			c.Locals("sessionID", sessionId)
			return c.Next()
		} else if sessionId == "" {
			c.Status(fiber.ErrForbidden.Code)
			return c.JSON(response.BadRes(ErrSessionRequired))
		}
		return fiber.ErrUpgradeRequired
	}
}

func (mc *MonitorController) GetEventsCount() fiber.Handler {
	return websocket.New(func(c *websocket.Conn) {
		for {
			_ = c.Locals("sessionID").(string)
		}
	})
}
