package reqs

import (
	"github.com/Izumra/SKUD_OKEI/domain/dto/integrserv"
)

type ReqEventFilter struct {
	BeginTime  string
	EndTime    string
	EventTypes []*integrserv.EventType
	Persons    []*integrserv.PersonData
}
