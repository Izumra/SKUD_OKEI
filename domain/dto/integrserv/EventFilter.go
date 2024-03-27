package integrserv

import (
	"encoding/xml"
	"time"
)

type EventFilter struct {
	XMLName    xml.Name
	BeginTime  time.Time
	EndTime    time.Time
	EventTypes []EventTypes
	Persons    []PersonData
}
