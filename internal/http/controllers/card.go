package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/Izumra/SKUD_OKEI/domain/dto/integrserv"
	"github.com/Izumra/SKUD_OKEI/domain/dto/reqs"
	"github.com/Izumra/SKUD_OKEI/internal/lib/response"
	"github.com/gofiber/fiber/v2"
)

type CardService interface {
	GetKeys(ctx context.Context, sessionId string, offset int64, count int64) ([]*integrserv.KeyData, error)
	GetKeyData(ctx context.Context, sessionId string, cardNo string) (*integrserv.KeyData, error)
	UpdateKeyData(ctx context.Context, sessionId string, keyData *integrserv.KeyData) (*integrserv.KeyData, error)
	AddKey(ctx context.Context, sessionId string, keyData *integrserv.KeyData) (*integrserv.KeyData, error)
	ReadKeyCode(ctx context.Context, sessionId string, idReader int) (string, error)
	ConvertWiegandToTouchMemory(ctx context.Context, sessionId string, code int, codeSize int) (string, error)
	ConvertPinToTouchMemory(ctx context.Context, sessionId string, pin string) (string, error)
}

type CardController struct {
	service CardService
}

func RegistrCardAPI(router fiber.Router, cs CardService) {
	ac := CardController{
		service: cs,
	}

	router.Get("/by_card_number/:card_no", ac.GetKeyData())
	router.Put("/", ac.UpdateKeyData())
	router.Post("/", ac.AddKey())
	router.Get("/read_card_number/:id_reader", ac.ReadCardNumber())
	router.Post("/wiegand_to_touch_memory", ac.WiegandToTouchMemory())
	router.Post("/pin_to_touch_memory/:code", ac.PinToTouchMemory())
	router.Get("/:offset/:count", ac.GetKeys())

}

// @Summary Получение ключей СКУД
// @Description Метод API, позволяющий авторизированному пользователю получить список ключей начиная с шага смещения, указанного в параметре 'offset', количества, заданного параметром 'count'
// @Tags Keys
// @Produce  json
// @Param offset path int true "Шаг смещения" default(0)
// @Param count path int true "Количество" default(0)
// @Success 200 {object} response.Body{data=[]integrserv.KeyData,error=nil} "Структура успешного ответа запроса получения ключей"
// @Failure 404 {object} response.Body{data=nil} "Структура неудачного ответа запроса получения ключей"
// @Router /api/cards/ [get]
func (cc *CardController) GetKeys() fiber.Handler {
	return func(c *fiber.Ctx) error {
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
			return c.JSON(response.BadRes(fmt.Errorf("Неверный формат параметра количества ключей")))
		}

		result, err := cc.service.GetKeys(c.Context(), session, offset, count)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.JSON(response.BadRes(err))
		}

		return c.JSON(response.SuccessRes(result))
	}
}

// @Summary Получение данных о ключе СКУД
// @Description Метод API, позволяющий авторизированному пользователю получить информацию о ключе по переданному параметру электронного кода 'card_no'
// @Tags Keys
// @Produce  json
// @Param card_no path string true "Электронный код пропуска" default("CA00000082942101")
// @Success 200 {object} response.Body{data=integrserv.KeyData,error=nil} "Структура успешного ответа запроса получения данных о ключе"
// @Failure 404 {object} response.Body{data=nil} "Структура неудачного ответа запроса получения данных о ключе"
// @Router /api/cards/by_card_number [get]
func (cc *CardController) GetKeyData() fiber.Handler {
	return func(c *fiber.Ctx) error {
		session := c.Cookies("session", "")

		cardNumberParam := c.Params("card_no")

		if cardNumberParam == "" {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(response.BadRes(fmt.Errorf("Номер карты не может быть пустым")))
		}

		result, err := cc.service.GetKeyData(c.Context(), session, cardNumberParam)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.JSON(response.BadRes(err))
		}

		return c.JSON(response.SuccessRes(result))
	}
}

// @Summary Считывание электронного кода
// @Description Метод API, позволяющий авторизированному пользователю считать электронный код пропуска со считывателя
// @Tags Keys
// @Produce  json
// @Param id_reader path int true "Номер считывателя" default(2)
// @Success 200 {object} response.Body{data=integrserv.KeyData,error=nil} "Структура успешного ответа запроса получения электронного кода пропуска"
// @Failure 404 {object} response.Body{data=nil} "Структура неудачного ответа запроса получения электронного кода пропуска"
// @Router /api/cards/read_card_number [get]
func (cc *CardController) ReadCardNumber() fiber.Handler {
	return func(c *fiber.Ctx) error {
		session := c.Cookies("session", "")

		idReaderParam := c.Params("id_reader", "0")

		idReader, err := strconv.Atoi(idReaderParam)
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(response.BadRes(fmt.Errorf("Неверный формат идентификатора считывателя")))
		}

		result, err := cc.service.ReadKeyCode(c.Context(), session, idReader)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.JSON(response.BadRes(err))
		}

		return c.JSON(response.SuccessRes(result))
	}
}

// @Summary Регистрация электронного ключа
// @Description Метод API, позволяющий авторизированному пользователю зарегестрировать новый электронный ключ
// @Tags Keys
// @Accept json
// @Produce  json
// @Param NewKeyBody body integrserv.KeyData true "Тело запроса формата 'application/json', содержащее информацию о добавляемом ключе"
// @Success 200 {object} response.Body{data=integrserv.KeyData,error=nil} "Структура успешного ответа запроса добавления нового ключа"
// @Failure 500 {object} response.Body{data=nil} "Структура неудачного ответа запроса добавления нового ключа"
// @Router /api/cards [post]
func (cc *CardController) AddKey() fiber.Handler {
	return func(c *fiber.Ctx) error {
		session := c.Cookies("session")

		var body integrserv.KeyData

		if err := json.Unmarshal(c.Body(), &body); err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(response.BadRes(fmt.Errorf("Неверный формат данных добавляемого ключа")))
		}

		result, err := cc.service.AddKey(c.Context(), session, &body)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.JSON(response.BadRes(err))
		}

		return c.JSON(response.SuccessRes(result))
	}
}

// @Summary Обновление информации электронного ключа
// @Description Метод API, позволяющий авторизированному пользователю изменить информацию об электронном ключе
// @Tags Keys
// @Accept json
// @Produce  json
// @Param UpdateKeyBody body integrserv.KeyData true "Тело запроса формата 'application/json', содержащее информацию для обновления ключа"
// @Success 200 {object} response.Body{data=integrserv.KeyData,error=nil} "Структура успешного ответа запроса обновления ключа"
// @Failure 404 {object} response.Body{data=nil} "Структура неудачного ответа запроса обновления ключа"
// @Router /api/cards [put]
func (cc *CardController) UpdateKeyData() fiber.Handler {
	return func(c *fiber.Ctx) error {
		session := c.Cookies("session")

		var body integrserv.KeyData

		if err := json.Unmarshal(c.Body(), &body); err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(response.BadRes(fmt.Errorf("Неверный формат данных изменяемого ключа")))
		}

		result, err := cc.service.UpdateKeyData(c.Context(), session, &body)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.JSON(response.BadRes(err))
		}

		return c.JSON(response.SuccessRes(result))
	}
}

// @Summary Конвертация десятичного кода карточки к формату Wiegand 26
// @Description Метод API, позволяющий авторизированному пользователю конвертировать десятичный код на личецвой стороне электронного пропуска к формату, понятному 'Орион Про' - Wiegand 26
// @Tags Keys
// @Accept json
// @Produce  json
// @Param ConvertKeyCodeToWiegandBody body reqs.WiegandToTouchMemory true "Тело запроса формата 'application/json', содержащее информацию для конвертации электронного кода"
// @Success 200 {object} response.Body{data=string,error=nil} "Структура успешного ответа запроса конвертации кода ключа"
// @Failure 404 {object} response.Body{data=nil} "Структура неудачного ответа запроса конвертации кода ключа"
// @Router /api/cards/wiegand_to_touch_memory [post]
func (cc *CardController) WiegandToTouchMemory() fiber.Handler {
	return func(c *fiber.Ctx) error {
		session := c.Cookies("session")

		var body reqs.WiegandToTouchMemory
		if err := json.Unmarshal(c.Body(), &body); err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(response.BadRes(ErrBodyParse))
		}

		code, err := cc.service.ConvertWiegandToTouchMemory(c.Context(), session, body.Code, body.CodeSize)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.JSON(response.BadRes(err))
		}

		return c.JSON(response.SuccessRes(code))
	}
}

func (cc *CardController) PinToTouchMemory() fiber.Handler {
	return func(c *fiber.Ctx) error {
		session := c.Cookies("session")
		pinCode := c.Params("code")

		code, err := cc.service.ConvertPinToTouchMemory(c.Context(), session, pinCode)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.JSON(response.BadRes(err))
		}

		return c.JSON(response.SuccessRes(code))
	}
}
