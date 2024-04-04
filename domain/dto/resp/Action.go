package resp

import "time"

type Action struct {
	Time   time.Time
	Action string `value:"coming|leaving"`
}
