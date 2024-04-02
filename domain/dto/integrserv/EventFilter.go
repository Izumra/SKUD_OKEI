package integrserv

import (
	"encoding/xml"
	"time"
)

type EventFilter struct {
	XMLName     xml.Name
	BeginTime   time.Time
	EndTime     time.Time
	EventTypes  EventTypes
	Persons     Persons
	EntryPoints EntryPoints
}

type EventTypes struct {
	EventType []*EventType
}

type Persons struct {
	PersonData []*PersonData
}

type EntryPoints struct {
	EntryPoint []*EntryPoint
}

type EntryPoint struct {
	Id                int64
	Name              string
	EnterAccessZoneId int
	ExitAccessZoneId  int
}

type EventType struct {
	Id int64
}
