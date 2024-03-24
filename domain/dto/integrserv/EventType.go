package integrserv

type EventType struct {
	Id          int64
	CharId      string
	Description string
	Comments    string
	Category    string
	HexCode     string
	IsAlarm     bool
}
