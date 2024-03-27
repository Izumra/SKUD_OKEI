package reqs

import (
	"time"

	"github.com/Izumra/SKUD_OKEI/domain/dto/integrserv"
)

type EventFilter struct {
	BeginTime  time.Time
	EndTime    time.Time
	EventTypes []integrserv.EventTypes
	Persons    []integrserv.PersonData
}
