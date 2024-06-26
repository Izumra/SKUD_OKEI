package controllers

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	_ "github.com/Izumra/SKUD_OKEI/docs"
	"github.com/Izumra/SKUD_OKEI/domain/dto/reqs"
	"github.com/Izumra/SKUD_OKEI/domain/dto/resp"
	"github.com/Izumra/SKUD_OKEI/internal/lib/response"
	"github.com/Izumra/SKUD_OKEI/internal/services/auth"
	"github.com/gofiber/fiber/v2"
)

var (
	ErrSessionNotFound = errors.New("Сессия пользователя не была найдена")
)

type AuthService interface {
	Login(ctx context.Context, username, password string) (*resp.SuccessAuth, error)
	Registrate(ctx context.Context, username, password string) (*resp.SuccessAuth, error)
}

type AuthController struct {
	sessionStorage auth.SessionStorage
	service        AuthService
}

func RegistrAuthAPI(router fiber.Router, as AuthService, ss auth.SessionStorage) {
	ac := AuthController{
		sessionStorage: ss,
		service:        as,
	}

	router.Post("/login", ac.Login)
	router.Post("/logout", ac.Logout)
	router.Post("/registrate", ac.Registrate)
}

// @Summary Завершение сессии
// @Description Метод API для принудительного завершения сессии пользователя
// @Tags Auth
// @Produce  json
// @Success 200 {object} response.Body{data=string,error=nil} "Завершение сессии"
// @Router /logout [post]
func (ac *AuthController) Logout(c *fiber.Ctx) error {
	c.ClearCookie("session")
	return c.JSON(response.SuccessRes("Пользователь вышел"))
}

// @Summary Авторизация
// @Description Метод API для предоставления доступа пользователю к системе
// @Tags Auth
// @Accept  json
// @Produce  json
// @Param LoginBody body reqs.LoginBody true "Тело запроса авторизации пользователя формата 'application/json', в котором передается имя пользователя и пароль"
// @Success 200 {object} response.Body{data=resp.SuccessAuth,error=nil} "Успешная авторизация"
// @Failure 404 {object} response.Body{data=nil} "Авторизация с недествительными данными"
// @Router /login [post]
func (ac *AuthController) Login(c *fiber.Ctx) error {
	sessionId := c.Cookies("session", "")
	if sessionId != "" {
		session, err := ac.sessionStorage.GetByID(c.Context(), sessionId)
		if err != nil {
			var data reqs.LoginBody
			err := json.Unmarshal(c.Body(), &data)
			if err != nil {
				c.Status(fiber.StatusBadRequest)
				return c.JSON(response.BadRes(ErrBodyParse))
			}

			result, err := ac.service.Login(c.Context(), data.Username, data.Password)
			if err != nil {
				c.Status(fiber.StatusInternalServerError)
				return c.JSON(response.BadRes(err))
			}

			c.Cookie(&fiber.Cookie{
				Name:     "session",
				Value:    result.SessionId,
				MaxAge:   48 * 60 * 60 * 1000,
				SameSite: "Strict",
				Expires:  time.Now().Add(48 * time.Hour),
			})

			return c.JSON(response.SuccessRes(result))
		}

		return c.JSON(response.SuccessRes(resp.SuccessAuth{
			Username: session.Username,
		}))
	}

	var data reqs.LoginBody
	err := json.Unmarshal(c.Body(), &data)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(response.BadRes(ErrBodyParse))
	}

	result, err := ac.service.Login(c.Context(), data.Username, data.Password)
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(response.BadRes(err))
	}

	c.Cookie(&fiber.Cookie{
		Name:     "session",
		Value:    result.SessionId,
		MaxAge:   48 * 60 * 60 * 1000,
		SameSite: "Strict",
		Expires:  time.Now().Add(48 * time.Hour),
	})

	return c.JSON(response.SuccessRes(result))
}

func (ac *AuthController) Registrate(c *fiber.Ctx) error {
	var data reqs.RegBody
	err := json.Unmarshal(c.Body(), &data)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(response.BadRes(ErrBodyParse))
	}

	result, err := ac.service.Registrate(c.Context(), data.Username, data.Password)
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(response.BadRes(err))
	}

	c.Cookie(&fiber.Cookie{
		Name:     "session",
		Value:    result.SessionId,
		MaxAge:   48 * 60 * 60 * 1000,
		SameSite: "Strict",
		Expires:  time.Now().Add(48 * time.Hour),
	})
	return c.JSON(response.SuccessRes(result))
}
