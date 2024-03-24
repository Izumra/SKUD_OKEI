package http

import (
	"github.com/Izumra/SKUD_OKEI/internal/http/controllers"
	"github.com/gofiber/fiber/v3"
)

func RegistrHandlers(
	app *fiber.App,
	authService controllers.AuthService,
	personService controllers.PersonsService,
) {
	api := app.Group("/api")

	controllers.RegistrAuthAPI(app, authService)

	personsRouter := api.Group("/persons")
	controllers.RegistrPersonsAPI(personsRouter, personService)
}
