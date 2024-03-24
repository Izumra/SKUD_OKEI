package entity

import valueobject "github.com/Izumra/SKUD_OKEI/domain/value-object"

type User struct {
	Id       int64
	Username string
	Password string
	Role     valueobject.Role
}
