package integrserv

import "time"

const (
	EventItemSECTION = iota
	EventItemLOOP
	EventItemDEVICE
	EventItemREADER
	EventItemRELAY
	EventItemACCESSZONE
	EventItemACCESSPOINT
	EventItemSECTIONGROUP
)

type Event struct {
	EventId       string
	EventDate     time.Time
	PassMode      int
	LastName      string
	FirstName     string
	MiddleName    string
	TabNum        string
	PersonId      int64
	CardNo        string
	Description   string
	AccessPointId int
}
