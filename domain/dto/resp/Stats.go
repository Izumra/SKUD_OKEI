package resp

import "github.com/Izumra/SKUD_OKEI/domain/dto/integrserv"

type Stats struct {
	CountInside  int
	CountOutside int
	AnomalyIn    int
	AnomalyOut   int
	Events       []integrserv.Event
}
