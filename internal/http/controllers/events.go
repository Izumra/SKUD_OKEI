package controllers

import (
	"context"
	"encoding/json"

	"github.com/Izumra/SKUD_OKEI/domain/dto/reqs"
	"github.com/Izumra/SKUD_OKEI/internal/lib/resp"
	"github.com/gofiber/fiber/v3"
)

type EventsService interface {
	Login(ctx context.Context, username, password string) (string, error)
	Registrate(ctx context.Context, username, password string) (string, error)
}

type EventsController struct {
	service EventsService
}

func RegistrEventAPI(router fiber.Router, es EventsService) {
	ec := EventsController{
		service: es,
	}

	router.Post("/count", ec.GetEventsCount())
	router.Post("/", ec.GetEvents())
}

func (ec *EventsController) GetEventsCount() fiber.Handler {
	return func(c fiber.Ctx) error {
		var data reqs.LoginBody
		err := json.Unmarshal(c.Body(), &data)
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(resp.BadRes(ErrBodyParse))
		}

		result, err := ec.service.Login(c.Context(), data.Username, data.Password)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.JSON(resp.BadRes(err))
		}

		return c.JSON(resp.SuccessRes(result))
	}
}

func (ec *EventsController) GetEvents() fiber.Handler {
	return func(c fiber.Ctx) error {
		var data reqs.RegBody
		err := json.Unmarshal(c.Body(), &data)
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(resp.BadRes(ErrBodyParse))
		}

		result, err := ec.service.Registrate(c.Context(), data.Username, data.Password)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.JSON(resp.BadRes(err))
		}

		return c.JSON(resp.SuccessRes(result))
	}
}
