package integrserv

import (
	"time"
)

type PersonData struct {
	Id         int64
	FirstName  string
	LastName   string
	MiddleName string
	Status     int
	ChangeTime time.Time
}
