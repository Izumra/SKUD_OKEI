package controllers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
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

	router.Get("/count", pc.GetPersonsCount)
	router.Get("/departments", pc.GetDepartments)
	router.Post("/", pc.AddPerson)
	router.Delete("/", pc.DeletePerson)
	router.Put("/", pc.UpdatePerson)
	router.Post("/filter/:offset/:count/", pc.GetPersons)
	router.Get("/:id", pc.GetPersonById)
	router.Get("/activity/dayly/:date/:id", pc.GetDaylyUserStats)
	router.Get("/activity/monthly/:date/:id", pc.GetMonthlyUserStats)
}

// @Summary Фильтрация субъектов доступа СКУД
// @Description Метод API, позволяющий авторизированному пользователю получить список субъектов доступа по переданному в теле запроса, массиву фильтров формата 'ключ=значение'
// @Tags Persons
// @Accept  json
// @Produce  json
// @Param offset path int true "Шаг смещения" default(0)
// @Param count path int true "Количество" default(0)
// @Param Filters body []string true "Тело запроса формата 'application/json', содержащее массив фильтров"
// @Success 200 {object} response.Body{data=[]integrserv.PersonData,error=nil} "Структура успешного ответа запроса фильтрации субъектов"
// @Failure 404 {object} response.Body{data=nil} "Структура неудачного ответа запроса фильтрации субъектов"
// @Router /api/persons/filter/ [post]
func (pc *PersonsController) GetPersons(c *fiber.Ctx) error {
	session := c.Cookies("session", "")

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

func (pc *PersonsController) GetPersonsCount(c *fiber.Ctx) error {
	session := c.Cookies("session", "")

	result, err := pc.service.GetPersonsCount(c.Context(), session)
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(response.BadRes(err))
	}

	return c.JSON(response.SuccessRes(result))
}

// @Summary Получение информации о субъекте доступа СКУД
// @Description Метод API, позволяющий авторизированному пользователю получить информацию о субъекте доступа СКУД'
// @Tags Persons
// @Produce  json
// @Success 200 {object} response.Body{data=integrserv.PersonData,error=nil} "Структура успешного ответа выполнения запроса получения информации о субъекте доступа СКУД"
// @Failure 404 {object} response.Body{data=nil} "Структура неудачного ответа выполнения запроса получения информации о субъекте доступа СКУД"
// @Router /api/persons/ [get]
func (pc *PersonsController) GetPersonById(c *fiber.Ctx) error {
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

// @Summary Добавление субъекта доступа СКУД
// @Description Метод API, позволяющий авторизированному пользователю добавить нового субъекта доступа СКУД
// @Tags Persons
// @Accept json
// @Produce json
// @Param AddPersonBody body integrserv.PersonData true "Тело запроса добавления нового субъекта доступа СКУД, формата 'application/json', в котором передается информация о субъекте доступа"
// @Success 200 {object} response.Body{data=integrserv.PersonData,error=nil} "Структура успешного ответа выполнения запроса добавления субъекта доступа СКУД"
// @Failure 400 {object} response.Body{data=nil} "Структура неудачного ответа выполнения запроса добавления субъекта доступа СКУД"
// @Router /api/persons/ [post]
func (pc *PersonsController) AddPerson(c *fiber.Ctx) error {
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

// @Summary Обновление информации о субъекте доступа СКУД
// @Description Метод API, позволяющий авторизированному пользователю изменить информацию о субъекте доступа СКУД
// @Tags Persons
// @Accept json
// @Produce json
// @Param UpdatePersonBody body integrserv.PersonData true "Тело запроса обновления информации о субъекте доступа СКУД, формата 'application/json', в котором передается информация о субъекте доступа"
// @Success 200 {object} response.Body{data=integrserv.PersonData,error=nil} "Структура успешного ответа выполнения запроса обновления информации о субъекте доступа СКУД"
// @Failure 400 {object} response.Body{data=nil} "Структура неудачного ответа выполнения запроса обновления информации о субъекте доступа СКУД"
// @Router /api/persons/ [put]
func (pc *PersonsController) UpdatePerson(c *fiber.Ctx) error {
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

// @Summary Удаление субъекта доступа СКУД
// @Description Метод API, позволяющий авторизированному пользователю удалить субъект доступа СКУД
// @Tags Persons
// @Accept json
// @Produce json
// @Param DeletePersonBody body integrserv.PersonData true "Тело запроса удаления субъекта доступа СКУД, формата 'application/json', в котором передается информация о субъекте доступа"
// @Success 200 {object} response.Body{data=integrserv.PersonData,error=nil} "Структура успешного ответа выполнения запроса удаления субъекта доступа СКУД"
// @Failure 400 {object} response.Body{data=nil} "Структура неудачного ответа выполнения запроса удаления субъекта доступа СКУД"
// @Router /api/persons/ [delete]
func (pc *PersonsController) DeletePerson(c *fiber.Ctx) error {
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

// @Summary Получение информации о группах субъектов доступа СКУД
// @Description Метод API, позволяющий авторизированному пользователю получить информацию о группах субъектов доступа СКУД
// @Tags Persons
// @Produce json
// @Success 200 {object} response.Body{data=[]integrserv.Department,error=nil} "Структура успешного ответа выполнения запроса получения информации о группах субъектов доступа СКУД"
// @Failure 500 {object} response.Body{data=nil} "Структура неудачного ответа выполнения запроса получения информации о группах субъектов доступа СКУД"
// @Router /api/persons/departments [get]
func (pc *PersonsController) GetDepartments(c *fiber.Ctx) error {
	session := c.Cookies("session", "")

	result, err := pc.service.GetDepartments(c.Context(), session)
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(response.BadRes(err))
	}

	return c.JSON(response.SuccessRes(result))
}

// @Summary Получение информации о статистике посещаемости субъектом доступа за день СКУД
// @Description Метод API, позволяющий авторизированному пользователю получить информацию о статистике посещаемости субъектом доступа за конкретный день
// @Tags Persons
// @Produce json
// @Param date path string true "День посещения" default(2024-03-02T15:04:05-07:00)
// @Param id path int true "Идентификатор субъекта доступа СКУД" default(1417)
// @Success 200 {object} response.Body{data=[]resp.Action,error=nil} "Структура успешного ответа выполнения запроса получения информации о статистике посещаемости субъектом доступа СКУД за конкретный день"
// @Failure 500 {object} response.Body{data=nil} "Структура неудачного ответа выполнения запроса получения информации о статистике посещаемости субъектом доступа СКУД за конкретный день"
// @Router /api/persons/activity/dayly [get]
func (pc *PersonsController) GetDaylyUserStats(c *fiber.Ctx) error {
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

// @Summary Получение информации о статистике посещаемости субъектом доступа за месяц СКУД
// @Description Метод API, позволяющий авторизированному пользователю получить информацию о статистике посещаемости субъектом доступа за конкретный месяц
// @Tags Persons
// @Produce json
// @Param date path string true "Месяц на который нужно расчитать статистику" default(2024-03-02T15:04:05-07:00)
// @Param id path int true "Идентификатор субъекта доступа СКУД" default(1417)
// @Success 200 {object} response.Body{data=[]resp.Activity,error=nil} "Структура успешного ответа выполнения запроса получения информации о статистике посещаемости субъектом доступа СКУД за конкретный месяц"
// @Failure 500 {object} response.Body{data=nil} "Структура неудачного ответа выполнения запроса получения информации о статистике посещаемости субъектом доступа СКУД за конкретный месяц"
// @Router /api/persons/activity/monthly [get]
func (pc *PersonsController) GetMonthlyUserStats(c *fiber.Ctx) error {
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
