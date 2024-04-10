package controllers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/Izumra/SKUD_OKEI/domain/dto/integrserv"
	"github.com/Izumra/SKUD_OKEI/domain/dto/resp"
	"github.com/Izumra/SKUD_OKEI/internal/lib/response"
	"github.com/gofiber/fiber/v2"
)

var (
	ErrBodyParse  = errors.New("Неверный формат тела запроса")
	ErrParamParse = errors.New("Неверный формат параметра запроса")
)

type PersonsService interface {
	GetPersons(ctx context.Context, sessionId string, offset int64, count int64, filters []string) ([]*integrserv.PersonData, error)
	GetPersonsCount(ctx context.Context, sessionId string) (int64, error)
	GetPersonById(ctx context.Context, sessionId string, id int64) (*integrserv.PersonData, error)
	AddPerson(ctx context.Context, sessionId string, data integrserv.PersonData) (*integrserv.PersonData, error)
	UpdatePerson(ctx context.Context, sessionId string, data integrserv.PersonData) (*integrserv.PersonData, error)
	DeletePerson(ctx context.Context, sessionId string, data integrserv.PersonData) (*integrserv.PersonData, error)
	GetDepartments(ctx context.Context, sessionId string) ([]*integrserv.Department, error)
	GetDaylyUserStats(ctx context.Context, sessionId string, id int64, date time.Time) ([]*resp.Action, error)
	GetMonthlyUserStats(ctx context.Context, sessionId string, id int64, monthTime time.Time) ([]*resp.Activity, error)
}

type PersonsController struct {
	service PersonsService
}

func RegistrPersonsAPI(router fiber.Router, ps PersonsService) {
	pc := PersonsController{
		service: ps,
	}

	router.Get("/count", pc.GetPersonsCount())
	router.Get("/departments", pc.GetDepartments())
	router.Post("/", pc.AddPerson())
	router.Delete("/", pc.DeletePerson())
	router.Put("/", pc.UpdatePerson())
	router.Post("/filter/:offset/:count/", pc.GetPersons())
	router.Get("/:id", pc.GetPersonById())
	router.Get("/activity/dayly/:date/:id", pc.GetDaylyUserStats())
	router.Get("/activity/monthly/:date/:id", pc.GetMonthlyUserStats())
}

func (pc *PersonsController) GetPersons() fiber.Handler {
	return func(c *fiber.Ctx) error {
		session := c.Cookies("session", "")
		log.Println(session)

		offsetParam := c.Params("offset", "0")
		countParam := c.Params("count", "0")

		offset, err := strconv.ParseInt(offsetParam, 10, 0)
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(response.BadRes(fmt.Errorf("Неверный формат параметра смещения")))
		}

		count, err := strconv.ParseInt(countParam, 10, 0)
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(response.BadRes(fmt.Errorf("Неверный формат параметра количества пользователей")))
		}

		var body []string
		if err := json.Unmarshal(c.Body(), &body); err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(response.BadRes(fmt.Errorf("Неверный формат фильтров для запроса")))
		}

		result, err := pc.service.GetPersons(c.Context(), session, offset, count, body)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.JSON(response.BadRes(err))
		}

		return c.JSON(response.SuccessRes(result))
	}
}

func (pc *PersonsController) GetPersonsCount() fiber.Handler {
	return func(c *fiber.Ctx) error {
		session := c.Cookies("session", "")

		result, err := pc.service.GetPersonsCount(c.Context(), session)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.JSON(response.BadRes(err))
		}

		return c.JSON(response.SuccessRes(result))
	}
}

func (pc *PersonsController) GetPersonById() fiber.Handler {
	return func(c *fiber.Ctx) error {
		session := c.Cookies("session", "")

		idParam := c.Params("id", "0")

		id, err := strconv.ParseInt(idParam, 10, 0)
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(response.BadRes(fmt.Errorf("Неверный формат id пользователя")))
		}

		result, err := pc.service.GetPersonById(c.Context(), session, id)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.JSON(response.BadRes(err))
		}

		return c.JSON(response.SuccessRes(result))
	}
}

func (pc *PersonsController) AddPerson() fiber.Handler {
	return func(c *fiber.Ctx) error {
		session := c.Cookies("session", "")

		var data integrserv.PersonData
		err := json.Unmarshal(c.Body(), &data)
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(response.BadRes(ErrBodyParse))
		}

		result, err := pc.service.AddPerson(c.Context(), session, data)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.JSON(response.BadRes(err))
		}

		return c.JSON(response.SuccessRes(result))
	}
}

func (pc *PersonsController) UpdatePerson() fiber.Handler {
	return func(c *fiber.Ctx) error {
		session := c.Cookies("session", "")

		var data integrserv.PersonData
		err := json.Unmarshal(c.Body(), &data)
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(response.BadRes(ErrBodyParse))
		}

		result, err := pc.service.UpdatePerson(c.Context(), session, data)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.JSON(response.BadRes(err))
		}

		return c.JSON(response.SuccessRes(result))
	}
}

func (pc *PersonsController) DeletePerson() fiber.Handler {
	return func(c *fiber.Ctx) error {
		session := c.Cookies("session", "")

		var data integrserv.PersonData
		err := json.Unmarshal(c.Body(), &data)
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(response.BadRes(ErrBodyParse))
		}

		result, err := pc.service.DeletePerson(c.Context(), session, data)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.JSON(response.BadRes(err))
		}

		return c.JSON(response.SuccessRes(result))
	}
}

func (pc *PersonsController) GetDepartments() fiber.Handler {
	return func(c *fiber.Ctx) error {

		session := c.Cookies("session", "")

		result, err := pc.service.GetDepartments(c.Context(), session)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.JSON(response.BadRes(err))
		}

		return c.JSON(response.SuccessRes(result))
	}
}

func (pc *PersonsController) GetDaylyUserStats() fiber.Handler {
	return func(c *fiber.Ctx) error {
		session := c.Cookies("session", "")

		idParam := c.Params("id", "0")

		id, err := strconv.ParseInt(idParam, 10, 0)
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(response.BadRes(fmt.Errorf("Неверный формат id пользователя")))
		}

		layout := "2006-01-02T15:04:05-07:00"
		date, err := time.Parse(layout, c.Params("date", time.Now().Format(layout)))
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(response.BadRes(ErrParamParse))
		}

		result, err := pc.service.GetDaylyUserStats(c.Context(), session, id, date)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.JSON(response.BadRes(err))
		}

		return c.JSON(response.SuccessRes(result))
	}
}

func (pc *PersonsController) GetMonthlyUserStats() fiber.Handler {
	return func(c *fiber.Ctx) error {
		session := c.Cookies("session", "")

		idParam := c.Params("id", "0")

		id, err := strconv.ParseInt(idParam, 10, 0)
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(response.BadRes(fmt.Errorf("Неверный формат id пользователя")))
		}

		layout := "2006-01-02T15:04:05-07:00"
		monthTime, err := time.Parse(layout, c.Params("date", time.Now().Format(layout)))
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(response.BadRes(ErrParamParse))
		}

		result, err := pc.service.GetMonthlyUserStats(c.Context(), session, id, monthTime)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.JSON(response.BadRes(err))
		}

		return c.JSON(response.SuccessRes(result))
	}
}
