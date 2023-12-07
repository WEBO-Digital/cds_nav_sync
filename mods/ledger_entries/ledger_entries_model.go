package ledgerentries

// Create represents the Create element.
type LedgerEntriesCreate struct {
	CurrentJnlBatchName string        `xml:"CurrentJnlBatchName" json:"currentJnlBatchName"`
	VendorPayment       VendorPayment `xml:"VendorPayment" json:"vendorPayment"`
}

// VendorPayment represents the VendorPayment element.
type VendorPayment struct {
	PostingDate      string  `xml:"Posting_Date" json:"postingDate"`
	DocumentDate     string  `xml:"Document_Date" json:"documentDate"`
	DocumentType     string  `xml:"Document_Type" json:"documentType"`
	AccountType      string  `xml:"Account_Type" json:"accountType"`
	AccountNo        string  `xml:"Account_No" json:"accountNo"`
	Amount           float64 `xml:"Amount" json:"amount"`
	AppliesToDocType string  `xml:"Applies_to_Doc_Type" json:"appliesToDocType"`
	AppliesToDocNo   string  `xml:"Applies_to_Doc_No" json:"appliesToDocNo"`
}

type BackToCDSLedgerEntriesResponse struct {
	DocumentNo int `json:"document_no,omitempty"`
}
