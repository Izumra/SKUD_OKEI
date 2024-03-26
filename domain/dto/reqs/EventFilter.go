package reqs

import (
	"encoding/xml"
	"time"

	"github.com/Izumra/SKUD_OKEI/domain/dto/integrserv"
)

type EventFilter struct {
	XMLName xml.Name `json:"-"`

	BeginTime  time.Time
	EndTime    time.Time
	EventTypes []integrserv.EventTypes
	Persons    []integrserv.PersonData
}
