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

// <SOAP-ENV:TEvent id="36" xsi:type="NS2:TEvent">
//
//	    <EventId xsi:type="xsd:string">{62BCF9B1-06C0-EE11-9690-A87EEA245DFF}</EventId>
//	    <EventDate xsi:type="xsd:dateTime">2024-01-31T12:02:30.000+05:00</EventDate>
//	    <Description xsi:type="xsd:string">2: Выход   Дверь 2,   Считыватель 2, Прибор 3</Description>
//	    <ZoneAddress xsi:type="xsd:int">2</ZoneAddress>
//	    <AccessPointId xsi:type="xsd:int">2</AccessPointId>
//	    <AccessZoneId xsi:type="xsd:int">1</AccessZoneId>
//	    <PassMode xsi:type="xsd:int">2</PassMode>
//	    <CardNo xsi:type="xsd:string">08003500673CEC01</CardNo>
//	    <PersonId xsi:type="xsd:int">2123</PersonId>
//	    <LastName xsi:type="xsd:string">Кузьмина </LastName>
//	    <FirstName xsi:type="xsd:string">Полина</FirstName>
//	    <MiddleName xsi:type="xsd:string">Сергеевна</MiddleName>
//	    <ItemId xsi:type="xsd:int">11</ItemId>
//	    <ItemType xsi:type="xsd:string">READER</ItemType>
//	</SOAP-ENV:TEvent>

// <SOAP-ENV:TEvent id="38" xsi:type="NS2:TEvent">
//
//	    <EventId xsi:type="xsd:string">{63BCF9B1-06C0-EE11-9690-A87EEA245DFF}</EventId>
//	    <EventDate xsi:type="xsd:dateTime">2024-01-31T12:02:36.000+05:00</EventDate>
//	    <Description xsi:type="xsd:string">2: Вход   Дверь 2,   Считыватель 1, Прибор 3</Description>
//	    <ZoneAddress xsi:type="xsd:int">1</ZoneAddress>
//	    <AccessPointId xsi:type="xsd:int">2</AccessPointId>
//	    <AccessZoneId xsi:type="xsd:int">2</AccessZoneId>
//	    <PassMode xsi:type="xsd:int">1</PassMode>
//	    <CardNo xsi:type="xsd:string">D700000021D6DE01</CardNo>
//	    <PersonId xsi:type="xsd:int">1873</PersonId>
//	    <LastName xsi:type="xsd:string">Совертков</LastName>
//	    <FirstName xsi:type="xsd:string">Семен</FirstName>
//	    <MiddleName xsi:type="xsd:string">Евгеньевич</MiddleName>
//	    <ItemId xsi:type="xsd:int">10</ItemId>
//	    <ItemType xsi:type="xsd:string">READER</ItemType>
//	</SOAP-ENV:TEvent>
type Event struct {
	EventId     string
	EventTypeId int64
	EventDate   time.Time
	Description string
	ItemId      int64
	ItemType    int
}
