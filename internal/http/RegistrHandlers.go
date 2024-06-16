package http

import (
	_ "github.com/Izumra/SKUD_OKEI/docs"
	"github.com/Izumra/SKUD_OKEI/internal/http/controllers"
	"github.com/Izumra/SKUD_OKEI/internal/http/controllers/ws"
	"github.com/Izumra/SKUD_OKEI/internal/services/auth"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/swagger"
)

// @title          API спецификация для СКУД ГАПОУ "ОКЭИ"
// @version        1.0
// @description    This is a sample fiber project server.
// @termsOfService http://swagger.io/terms/
// @contact.name Марк Загуменников
// @contact.email zagumennikovmark@gmail.com
// @host localhost:8082
// @BasePath /
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

	app.Get("/swagger/*", swagger.HandlerDefault)

	app.Static("/", "dist")

	api := app.Group("/api")

	controllers.RegistrAuthAPI(app, authService, sessionStorage)

	personsRouter := api.Group("/persons")
	controllers.RegistrPersonsAPI(personsRouter, personService)

	eventsRouter := api.Group("/events")
	controllers.RegistrEventAPI(eventsRouter, eventsService, sessionStorage)

	cardRouter := api.Group("/cards")
	controllers.RegistrCardAPI(cardRouter, cardService)

	webSocketRouter := api.Group("/ws")
	ws.RegistrWSAPI(webSocketRouter, eventsService, sessionStorage)
}
