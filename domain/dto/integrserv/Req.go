package integrserv

import "encoding/xml"

type EnvelopeReq struct {
	XMLName  xml.Name `xml:"soap:Envelope" json:"-"`
	XmlsSoap string   `xml:"xmlns:soap,attr" json:"-"`
	BodySoap BodyReq  `xml:"soap:Body" json:"-"`
}

type BodyReq struct {
	Data any
}
