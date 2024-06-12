package ws

import (
	"context"
	"encoding/xml"
	"errors"
	"time"

	"github.com/Izumra/SKUD_OKEI/domain/dto/integrserv"
	"github.com/Izumra/SKUD_OKEI/domain/dto/resp"
	"github.com/Izumra/SKUD_OKEI/internal/lib/response"
	"github.com/Izumra/SKUD_OKEI/internal/services/auth"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

var (
	ErrSessionRequired = errors.New("Для продолжения действия требуется сессия")
	ErrWrongReqBody    = errors.New("Тело сообщения неверного формата")
)

type WSService interface {
	GetEvents(ctx context.Context, eventsFilter *integrserv.EventFilter) ([]integrserv.Event, error)
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
		sessionId := c.Cookies("session", "")

		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("sessionID", sessionId)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	}
}

func (mc *WSController) Monitor() fiber.Handler {
	return websocket.New(func(c *websocket.Conn) {

		ctx, cancel := context.WithCancel(context.TODO())
		defer cancel()

		var lastUpdate time.Time

		now := time.Now()
		lastUpdate = time.Date(now.Year(), now.Month(), now.Day(), 6, 0, 0, 0, now.Location())

		recentlyRecords := map[string]bool{}

		var evnts []integrserv.Event
		users := make(map[int64]integrserv.Event)
		stats := &resp.Stats{}

		var closeHandlerSetted bool
		for {
			select {
			case <-ctx.Done():
				return
			default:
				if !closeHandlerSetted {
					go func() {
						_, _, err := c.ReadMessage()
						if err != nil {
							if websocket.IsUnexpectedCloseError(err) || websocket.IsCloseError(err) {
								cancel()
							}
						}
					}()
					closeHandlerSetted = true
				}

				if len(evnts) != 0 {
					time.Sleep(1 * time.Second)
				}

				filter := integrserv.EventFilter{
					XMLName: xml.Name{
						Local: "GetEvents",
					},
					BeginTime: lastUpdate,
					EndTime:   lastUpdate.Add(time.Hour),
					Offset:    0,
					Count:     100,
				}

				events, err := mc.service.GetEvents(ctx, &filter)
				if err != nil {
					c.WriteJSON(response.BadRes(err))
					return
				}

				if len(events) != 0 {
					if len(recentlyRecords) == 0 {
						for i := range events {
							recentlyRecords[events[i].EventId] = true

							UpdateStats(stats, events[i], users)

							users[events[i].PersonId] = events[i]
						}
						evnts = events
					} else {
						for i := range events {
							if _, ok := recentlyRecords[events[i].EventId]; !ok {
								recentlyRecords[events[i].EventId] = true
								evnts = append(evnts, events[i])

								UpdateStats(stats, events[i], users)

								users[events[i].PersonId] = events[i]

							}
						}
					}

					stats.Events = evnts
					lastUpdate = evnts[len(evnts)-1].EventDate

					err = c.Conn.WriteJSON(response.SuccessRes(stats))
					if err != nil {
						if websocket.IsUnexpectedCloseError(err) || websocket.IsCloseError(err) {
							cancel()
							break
						}
					}
				} else {
					lastUpdate = lastUpdate.Add(time.Hour)
				}
			}
		}
	})
}

func UpdateStats(stats *resp.Stats, event integrserv.Event, users map[int64]integrserv.Event) {
	if lastUsrEvent, ok := users[event.PersonId]; ok {
		if lastUsrEvent.PassMode == 1 && event.PassMode == 1 {
			stats.AnomalyIn++
		} else if lastUsrEvent.PassMode == 2 && event.PassMode == 2 {
			stats.AnomalyOut++
		}
	}
	if event.PassMode == 1 {
		stats.CountInside++
	} else if event.PassMode == 2 {
		stats.CountOutside++
	}
}
