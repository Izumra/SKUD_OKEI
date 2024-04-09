package ws

import (
	"context"
	"encoding/xml"
	"errors"
	"time"

	"github.com/Izumra/SKUD_OKEI/domain/dto/integrserv"
	"github.com/Izumra/SKUD_OKEI/internal/lib/response"
	"github.com/Izumra/SKUD_OKEI/internal/services/auth"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

var (
	ErrSessionRequired = errors.New(" Для продолжения действия требуется сессия")
	ErrWrongReqBody    = errors.New(" Тело сообщения неверного формата")
)

type WSService interface {
	GetEvents(ctx context.Context, eventsFilter *integrserv.EventFilter) ([]*integrserv.Event, error)
	GetEventsCount(ctx context.Context, eventsFilter *integrserv.EventCountFilter) (int64, error)
}

type WSController struct {
	sessStorage auth.SessionStorage
	service     WSService
}

func RegistrWSAPI(router fiber.Router, ws WSService, sessStorage auth.SessionStorage) {
	mc := WSController{
		sessStorage: sessStorage,
		service:     ws,
	}

	router.Use(mc.CheckRegisteredUpgrade())
	router.Get("/monitor", mc.Monitor())
}

func (mc *WSController) CheckRegisteredUpgrade() fiber.Handler {
	return func(c *fiber.Ctx) error {
		//sessionId := c.Get("Authorization")
		if websocket.IsWebSocketUpgrade(c) {
			//c.Locals("sessionID", sessionId)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	}
}

func (mc *WSController) Monitor() fiber.Handler {
	return websocket.New(func(c *websocket.Conn) {

		ctx, cancel := context.WithCancel(context.TODO())
		defer cancel()

		// session := c.Locals("sessionID").(string)

		// _, err := mc.sessStorage.GetByID(ctx, session)
		// if err != nil {
		// 	c.WriteJSON(response.BadRes(err))
		// 	return
		// }

		var lastUpdate time.Time

		now := time.Now()
		lastUpdate = time.Date(now.Year(), now.Month(), now.Day(), now.Hour()-1, 0, 0, 0, now.Location())

		recentlyRecords := map[string]bool{}

		type Stats struct {
			CountInside  int
			CountOutside int
			Events       []*integrserv.Event
		}

		var evnts []*integrserv.Event
		stats := Stats{}

		var closedHandlerSetted bool
		for {
			select {
			case <-ctx.Done():
				return
			default:
				if !closedHandlerSetted {
					go func() {
						_, _, err := c.ReadMessage()
						if err != nil {
							if websocket.IsUnexpectedCloseError(err) || websocket.IsCloseError(err) {
								cancel()
							}
						}
					}()
					closedHandlerSetted = true
				}

				if len(evnts) != 0 {
					time.Sleep(1 * time.Second)
				}

				filter := integrserv.EventFilter{
					XMLName: xml.Name{
						Local: "GetEvents",
					},
					BeginTime: lastUpdate,
					EndTime:   time.Now(),
					Offset:    0,
					Count:     100,
				}

				events, err := mc.service.GetEvents(ctx, &filter)
				if err != nil {
					c.WriteJSON(response.BadRes(err))
					cancel()
					return
				}

				if events != nil {
					eventsCount := len(evnts)

					var countNewRecords int
					if len(recentlyRecords) == 0 {
						for i := range events {
							recentlyRecords[events[i].EventId] = true
							if events[i].PassMode == 1 {
								stats.CountInside++
							} else if events[i].PassMode == 2 {
								stats.CountOutside++
							}
						}
						evnts = events
						countNewRecords = len(events)
					} else {
						for i := range events {
							if _, ok := recentlyRecords[events[i].EventId]; !ok {
								recentlyRecords[events[i].EventId] = true
								evnts = append(evnts, events[i])
								countNewRecords++

								if events[i].PassMode == 1 {
									stats.CountInside++
								} else if events[i].PassMode == 2 {
									stats.CountOutside++
								}
							}
						}
					}

					stats.Events = evnts[len(evnts)-countNewRecords:]
					lastUpdate = evnts[len(evnts)-1].EventDate
					if eventsCount < len(evnts) {
						err = c.Conn.WriteJSON(response.SuccessRes(stats))
						if err != nil {
							if websocket.IsUnexpectedCloseError(err) || websocket.IsCloseError(err) {
								cancel()
								break
							}
						}
					}
				}
			}
		}
	})
}
