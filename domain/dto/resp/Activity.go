package resp

import "time"

type Activity struct {
	Time    time.Time
	Coming  int
	Leaving int
}
