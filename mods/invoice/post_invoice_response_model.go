package invoice

import "encoding/xml"

// Save Data
type PostResponseInvoiceEnvelope struct {
	XMLName xml.Name                `xml:"Envelope"`
	Soap    string                  `xml:"xmlns:Soap,attr" json:"-"`
	Body    PostCreateResultInvoice `xml:"Body>PostPurchaseInvoice_Result" json:"body"`
}

type PostCreateResultInvoice struct {
	XMLName     xml.Name `xml:"PostPurchaseInvoice_Result"`
	XMLNS       string   `xml:"xmlns,attr" json:"-"`
	ReturnValue string   `xml:"return_value" json:"return_value"`
}
