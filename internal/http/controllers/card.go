package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/Izumra/SKUD_OKEI/domain/dto/integrserv"
	"github.com/Izumra/SKUD_OKEI/internal/lib/response"
	"github.com/gofiber/fiber/v2"
)

type CardService interface {
	GetKeys(ctx context.Context, sessionId string, offset int64, count int64) ([]*integrserv.KeyData, error)
	GetKeyData(ctx context.Context, sessionId string, cardNo string) (*integrserv.KeyData, error)
	UpdateKeyData(ctx context.Context, sessionId string, keyData *integrserv.KeyData) (*integrserv.KeyData, error)
	AddKey(ctx context.Context, sessionId string, keyData *integrserv.KeyData) (*integrserv.KeyData, error)
	ReadKeyCode(ctx context.Context, sessionId string, idReader int) (string, error)
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
	router.Get("/:offset/:count", ac.GetKeys())

}

func (cc *CardController) GetKeys() fiber.Handler {
	return func(c *fiber.Ctx) error {
		session := c.Cookies("session", "")

		offsetParam := c.Params("offset", "0")
		countParam := c.Params("count", "0")

		offset, err := strconv.ParseInt(offsetParam, 10, 0)
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(response.BadRes(fmt.Errorf(" Неверный формат параметра смещения")))
		}

		count, err := strconv.ParseInt(countParam, 10, 0)
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(response.BadRes(fmt.Errorf(" Неверный формат параметра количества ключей")))
		}

		result, err := cc.service.GetKeys(c.Context(), session, offset, count)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.JSON(response.BadRes(err))
		}

		return c.JSON(response.SuccessRes(result))
	}
}

func (cc *CardController) GetKeyData() fiber.Handler {
	return func(c *fiber.Ctx) error {
		session := c.Cookies("session", "")

		cardNumberParam := c.Params("card_no")

		if cardNumberParam == "" {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(response.BadRes(fmt.Errorf(" Номер карты не может быть пустым")))
		}

		result, err := cc.service.GetKeyData(c.Context(), session, cardNumberParam)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.JSON(response.BadRes(err))
		}

		return c.JSON(response.SuccessRes(result))
	}
}

func (cc *CardController) ReadCardNumber() fiber.Handler {
	return func(c *fiber.Ctx) error {
		session := c.Cookies("session", "")

		idReaderParam := c.Params("id_reader", "0")

		idReader, err := strconv.Atoi(idReaderParam)
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(response.BadRes(fmt.Errorf(" Неверный формат идентификатора считывателя")))
		}

		result, err := cc.service.ReadKeyCode(c.Context(), session, idReader)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.JSON(response.BadRes(err))
		}

		return c.JSON(response.SuccessRes(result))
	}
}

func (cc *CardController) AddKey() fiber.Handler {
	return func(c *fiber.Ctx) error {
		session := c.Cookies("session", "")

		var body integrserv.KeyData

		if err := json.Unmarshal(c.Body(), &body); err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(response.BadRes(fmt.Errorf(" Неверный формат данных добавляемого ключа")))
		}

		result, err := cc.service.UpdateKeyData(c.Context(), session, &body)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.JSON(response.BadRes(err))
		}

		return c.JSON(response.SuccessRes(result))
	}
}

func (cc *CardController) UpdateKeyData() fiber.Handler {
	return func(c *fiber.Ctx) error {
		session := c.Cookies("session", "")

		var body integrserv.KeyData

		if err := json.Unmarshal(c.Body(), &body); err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(response.BadRes(fmt.Errorf(" Неверный формат данных изменяемого ключа")))
		}

		result, err := cc.service.UpdateKeyData(c.Context(), session, &body)
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.JSON(response.BadRes(err))
		}

		return c.JSON(response.SuccessRes(result))
	}
}
