package integrserv

import "encoding/xml"

type EnvelopeResp struct {
	XMLName     xml.Name `xml:"Envelope" json:"-"`
	XmlsSoapEnv string   `xml:"SOAP-ENV,attr" json:"-"`
	XmlsSoapEnc string   `xml:"SOAP-ENC,attr" json:"-"`
	XmlsXsd     string   `xml:"xsd,attr" json:"-"`
	XmlsXsi     string   `xml:"xsi,attr" json:"-"`

	OperationResult any `xml:"Body"`
}
