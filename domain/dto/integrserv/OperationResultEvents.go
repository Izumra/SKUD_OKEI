package integrserv

type OperationResultEvents struct {
	SoapEnvEncodingStyle string `xml:"encodingStyle,attr" json:"-"`
	XmlnsNS1             string `xml:"NS1,attr" json:"-"`
	XmlnsNS2             string `xml:"NS2,attr" json:"-"`

	Result any `xml:"TEvent"`
}
