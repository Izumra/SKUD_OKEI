package integrserv

const (
	KeyCodeProgrammPass = iota
	KeyCodePinCode
	KeyCodePinTouchMemory
	KeyCodeProxy
	KeyCodeCar
	KeyCodeFingerPrint
	KeyCodeFacePrint
)

type KeyData struct {
	Id       int64
	CodeType int
	Code     string
	PersonId int64
}
