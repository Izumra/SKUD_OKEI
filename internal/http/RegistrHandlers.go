package http

import (
	"github.com/Izumra/SKUD_OKEI/internal/http/controllers"
	"github.com/Izumra/SKUD_OKEI/internal/http/controllers/ws"
	"github.com/Izumra/SKUD_OKEI/internal/services/auth"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func RegistrHandlers(
	app *fiber.App,
	sessionStorage auth.SessionStorage,
	authService controllers.AuthService,
	personService controllers.PersonsService,
	eventsService controllers.EventsService,
	cardService controllers.CardService,
) {
	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowOriginsFunc: func(origin string) bool { return true },
	}))

	api := app.Group("/api")

	controllers.RegistrAuthAPI(app, authService, sessionStorage)

	personsRouter := api.Group("/persons")
	controllers.RegistrPersonsAPI(personsRouter, personService)

	eventsRouter := api.Group("/events")
	controllers.RegistrEventAPI(eventsRouter, eventsService, sessionStorage)

	cardRouter := api.Group("/card")
	controllers.RegistrCardAPI(cardRouter, cardService)

	webSocketRouter := api.Group("/ws")
	ws.RegistrWSAPI(webSocketRouter, eventsService, sessionStorage)
}
