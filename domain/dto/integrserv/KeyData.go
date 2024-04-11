package integrserv

const (
	KeyCodeProgrammPass = iota + 1
	KeyCodePinCode
	KeyCodePinTouchMemory
	KeyCodeProxy
	KeyCodeCar
	KeyCodeFingerPrint
	KeyCodeFacePrint
)

type KeyData struct {
	Id              int64
	CodeType        int
	Code            string
	PersonId        int64
	AccessLevelId   int
	StartDate       string
	EndDate         string
	IsBlocked       bool
	IsStoreInDevice bool
	IsStoreInS2000  bool
	IsInStopList    bool
}
