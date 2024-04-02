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

// {
// 	"EventId": "{6B8D9F38-CEEC-EE11-9692-A87EEA245DFF}",
// 	"EventDate": "2024-03-28T11:41:36+05:00",
// 	"PassMode": 1,
// 	"LastName": "",
// 	"FirstName": "",
// 	"MiddleName": "",
// 	"TabNum": "",
// 	"PersonId": 0,
// 	"CardNo": "",
// 	"Description": "2: Вход  7800090082942101 Дверь 2,  7800090082942101 Считыватель ",
// 	"AccessPointId": 2
// }

// {
// 	"EventId": "{6B8D9F38-CEEC-EE11-9692-A87EEA245DFF}",
// 	"EventDate": "2024-03-28T11:41:36+05:00",
// 	"PassMode": 1,
// 	"LastName": "",
// 	"FirstName": "",
// 	"MiddleName": "",
// 	"TabNum": "",
// 	"PersonId": 0,
// 	"CardNo": "",
// 	"Description": "2: Вход  7800090082942101 Дверь 2,  7800090082942101 Считыватель ",
// 	"AccessPointId": 2
// },
// {
// 	"EventId": "{F9FBA73E-CEEC-EE11-9692-A87EEA245DFF}",
// 	"EventDate": "2024-03-28T11:41:46+05:00",
// 	"PassMode": 1,
// 	"LastName": "Загуменников ",
// 	"FirstName": "Марк ",
// 	"MiddleName": "Юрьевич",
// 	"TabNum": "",
// 	"PersonId": 764,
// 	"CardNo": "A300090091993501",
// 	"Description": "2: Вход   Дверь 2,   Считыватель 1, Прибор 3",
// 	"AccessPointId": 2
// },

// {
// 	"EventId": "{1FF5748F-82ED-EE11-9692-A87EEA245DFF}",
// 	"EventDate": "2024-03-29T09:12:32+05:00",
// 	"PassMode": 1,
// 	"LastName": "Кондауров ",
// 	"FirstName": "Владимир ",
// 	"MiddleName": "Максимович",
// 	"TabNum": "",
// 	"PersonId": 761,
// 	"CardNo": "BD000800B3D85801",
// 0011786328
// 	"Description": "1: Вход   Дверь 1,   Считыватель 1, Прибор 2",
// 	"AccessPointId": 1
// },

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
