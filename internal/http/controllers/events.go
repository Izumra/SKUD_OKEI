package controllers

import (
	"context"
	"encoding/xml"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/Izumra/SKUD_OKEI/domain/dto/integrserv"
	"github.com/Izumra/SKUD_OKEI/domain/dto/reqs"
	"github.com/Izumra/SKUD_OKEI/internal/lib/response"
	"github.com/Izumra/SKUD_OKEI/internal/services/auth"
	"github.com/gofiber/fiber/v2"
)

type EventsService interface {
	GetEvents(ctx context.Context, eventsFilter *integrserv.EventFilter) ([]*integrserv.Event, error)
	GetEventsCount(ctx context.Context, eventsFilter *integrserv.EventCountFilter) (int64, error)
}

type EventsController struct {
	sessStorage auth.SessionStorage
	service     EventsService
}

func RegistrEventAPI(router fiber.Router, es EventsService, ss auth.SessionStorage) {
	ec := EventsController{
		sessStorage: ss,
		service:     es,
	}

	router.Post("/count", ec.GetEventsCount())
	router.Post("/:offset/:count", ec.GetEvents())
}

func (ec *EventsController) GetEventsCount() fiber.Handler {
	return func(c *fiber.Ctx) error {

		_ = c.Cookies("session", "")

		reqBody := reqs.ReqEventFilter{}

		err := c.BodyParser(&reqBody)
		if err != nil {
			log.Println(err)
			c.Status(fiber.StatusBadRequest)
			return c.JSON(response.BadRes(ErrBodyParse))
		}

		layout := "2006-01-02T15:04:05"
		beginTime, err := time.Parse(layout, reqBody.BeginTime)
		if err != nil {
			log.Println(err)
			c.Status(fiber.StatusBadRequest)
			return c.JSON(response.BadRes(ErrBodyParse))
		}

		endTime, err := time.Parse(layout, reqBody.EndTime)
		if err != nil {
			log.Println(err)
			c.Status(fiber.StatusBadRequest)
			return c.JSON(response.BadRes(ErrBodyParse))
		}

		filter := integrserv.EventCountFilter{
			XMLName: xml.Name{
				Local: "GetEventsCount",
			},
			BeginTime: beginTime,
			EndTime:   endTime,
			EventTypes: integrserv.EventTypes{
				EventType: reqBody.EventTypes,
			},
			Persons: integrserv.Persons{
				PersonData: reqBody.Persons,
			},
		}
		result, err := ec.service.GetEventsCount(c.Context(), &filter)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.JSON(response.BadRes(err))
		}

		return c.JSON(response.SuccessRes(result))
	}
}

func (ec *EventsController) GetEvents() fiber.Handler {
	return func(c *fiber.Ctx) error {

		_ = c.Cookies("session", "")

		reqBody := reqs.ReqEventFilter{}

		err := c.BodyParser(&reqBody)
		if err != nil {
			log.Println(err)
			c.Status(fiber.StatusBadRequest)
			return c.JSON(response.BadRes(ErrBodyParse))
		}

		layout := "2006-01-02T15:04:05"
		beginTime, err := time.Parse(layout, reqBody.BeginTime)
		if err != nil {
			log.Println(err)
			c.Status(fiber.StatusBadRequest)
			return c.JSON(response.BadRes(ErrBodyParse))
		}

		endTime, err := time.Parse(layout, reqBody.EndTime)
		if err != nil {
			log.Println(err)
			c.Status(fiber.StatusBadRequest)
			return c.JSON(response.BadRes(ErrBodyParse))
		}

		offsetParam := c.Params("offset", "0")
		_, err = strconv.ParseInt(offsetParam, 10, 0)
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(response.BadRes(fmt.Errorf(" Неверный формат шага смещения")))
		}

		countParam := c.Params("count", "0")
		_, err = strconv.ParseInt(countParam, 10, 0)
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(response.BadRes(fmt.Errorf(" Неверный формат количества событий")))
		}

		filter := integrserv.EventFilter{
			XMLName: xml.Name{
				Local: "GetEvents",
			},
			BeginTime: beginTime,
			EndTime:   endTime,
			EventTypes: integrserv.EventTypes{
				EventType: reqBody.EventTypes,
			},
			Persons: integrserv.Persons{
				PersonData: reqBody.Persons,
			},
		}

		result, err := ec.service.GetEvents(c.Context(), &filter)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.JSON(response.BadRes(err))
		}

		return c.JSON(response.SuccessRes(result))
	}
}
