package controllers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/Izumra/SKUD_OKEI/domain/dto/integrserv"
	"github.com/Izumra/SKUD_OKEI/internal/lib/resp"
	"github.com/gofiber/fiber/v3"
)

var (
	ErrBodyParse = errors.New(" Неверный формат тела запроса")
)

type PersonsService interface {
	GetPersons(ctx context.Context, sessionId string, offset int64, count int64) (*integrserv.OperationResult, error)
	GetPersonsCount(ctx context.Context, sessionId string) (*integrserv.OperationResult, error)
	GetPersonById(ctx context.Context, sessionId string, id int64) (*integrserv.OperationResult, error)
	AddPerson(ctx context.Context, sessionId string, data integrserv.PersonData) (*integrserv.OperationResult, error)
	UpdatePerson(ctx context.Context, sessionId string, data integrserv.PersonData) (*integrserv.OperationResult, error)
	DeletePerson(ctx context.Context, sessionId string, data integrserv.PersonData) (*integrserv.OperationResult, error)
}

type PersonsController struct {
	service PersonsService
}

func RegistrPersonsAPI(router fiber.Router, ps PersonsService) {
	pc := PersonsController{
		service: ps,
	}

	router.Get("/count", pc.GetPersonsCount())
	router.Post("/", pc.AddPerson())
	router.Delete("/", pc.DeletePerson())
	router.Put("/", pc.UpdatePerson())
	router.Get("/:offset/:count", pc.GetPersons())
	router.Get("/:id", pc.GetPersonById())
}

func (pc *PersonsController) GetPersons() fiber.Handler {
	return func(c fiber.Ctx) error {
		session := c.Get("Authorization")

		offsetParam := c.Params("offset", "0")
		countParam := c.Params("count", "0")

		offset, err := strconv.ParseInt(offsetParam, 10, 0)
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(resp.BadRes(fmt.Errorf(" Неверный формат параметра смещения")))
		}

		count, err := strconv.ParseInt(countParam, 10, 0)
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(resp.BadRes(fmt.Errorf(" Неверный формат параметра количества пользователей")))
		}

		result, err := pc.service.GetPersons(c.Context(), session, offset, count)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.JSON(resp.BadRes(err))
		}

		return c.JSON(resp.SuccessRes(result))
	}
}

func (pc *PersonsController) GetPersonsCount() fiber.Handler {
	return func(c fiber.Ctx) error {
		session := c.Get("Authorization")

		result, err := pc.service.GetPersonsCount(c.Context(), session)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.JSON(resp.BadRes(err))
		}

		return c.JSON(resp.SuccessRes(result))
	}
}

func (pc *PersonsController) GetPersonById() fiber.Handler {
	return func(c fiber.Ctx) error {
		session := c.Get("Authorization")

		idParam := c.Params("id", "0")

		id, err := strconv.ParseInt(idParam, 10, 0)
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(resp.BadRes(fmt.Errorf(" Неверный формат id пользователя")))
		}

		result, err := pc.service.GetPersonById(c.Context(), session, id)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.JSON(resp.BadRes(err))
		}

		return c.JSON(resp.SuccessRes(result))
	}
}

func (pc *PersonsController) AddPerson() fiber.Handler {
	return func(c fiber.Ctx) error {
		session := c.Get("Authorization")

		var data integrserv.PersonData
		err := json.Unmarshal(c.Body(), &data)
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(resp.BadRes(ErrBodyParse))
		}

		result, err := pc.service.AddPerson(c.Context(), session, data)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.JSON(resp.BadRes(err))
		}

		return c.JSON(resp.SuccessRes(result))
	}
}

func (pc *PersonsController) UpdatePerson() fiber.Handler {
	return func(c fiber.Ctx) error {
		session := c.Get("Authorization")

		var data integrserv.PersonData
		err := json.Unmarshal(c.Body(), &data)
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(resp.BadRes(ErrBodyParse))
		}

		result, err := pc.service.UpdatePerson(c.Context(), session, data)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.JSON(resp.BadRes(err))
		}

		return c.JSON(resp.SuccessRes(result))
	}
}

func (pc *PersonsController) DeletePerson() fiber.Handler {
	return func(c fiber.Ctx) error {
		session := c.Get("Authorization")

		var data integrserv.PersonData
		err := json.Unmarshal(c.Body(), &data)
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(resp.BadRes(ErrBodyParse))
		}

		result, err := pc.service.DeletePerson(c.Context(), session, data)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.JSON(resp.BadRes(err))
		}

		return c.JSON(resp.SuccessRes(result))
	}
}
