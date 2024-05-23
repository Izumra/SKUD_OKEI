package integrserv

import (
	"encoding/json"
	"time"
)

const (
	KeyCodeProgrammPass = iota + 1
	KeyCodePinCode
	KeyCodePinTouchMemory
	KeyCodeProxy
	KeyCodeCar
	KeyCodeFingerPrint
	KeyCodeFacePrint
)

// 500009007FD4E701
// TODO: исправить временное решение с игнорированием полей при парсинге структуры json
type KeyData struct {
	Id              int64
	CodeType        int
	Code            string
	PersonId        int64
	AccessLevelId   int
	StartDate       time.Time
	EndDate         time.Time
	IsBlocked       bool
	IsStoreInDevice bool
	IsStoreInS2000  bool
	IsInStopList    bool
}

type jsonKeyData struct {
	StartDate string
	EndDate   string
	KeyData
}

func (kd *KeyData) UnmarshalJSON(data []byte) error {
	var jsonKeyData jsonKeyData

	if err := json.Unmarshal(data, &jsonKeyData); err != nil {
		return err
	}

	kd.Id = jsonKeyData.Id
	kd.CodeType = 4
	kd.Code = jsonKeyData.Code
	kd.PersonId = jsonKeyData.PersonId
	kd.AccessLevelId = 3
	kd.IsBlocked = false
	kd.IsStoreInDevice = true
	kd.IsStoreInS2000 = false
	kd.IsInStopList = false

	kd.StartDate = time.Now()
	kd.EndDate = time.Date(
		kd.StartDate.Year()+4,
		kd.StartDate.Month(),
		kd.StartDate.Day(),
		kd.StartDate.Hour(),
		kd.StartDate.Minute(),
		kd.StartDate.Second(),
		kd.StartDate.Nanosecond(),
		kd.StartDate.Location(),
	)

	return nil
}
