package ledgerentries

// Create represents the Create element.
type LedgerEntriesCreate struct {
	PaymentIDstring     int           `json:"payment_id,omitempty"`
	CurrentJnlBatchName string        `xml:"CurrentJnlBatchName" json:"currentJnlBatchName"`
	VendorPayment       VendorPayment `xml:"VendorPayment" json:"vendorPayment"`
}

// VendorPayment represents the VendorPayment element.
type VendorPayment struct {
	PostingDate      string  `xml:"Posting_Date,omitempty" json:"postingDate,omitempty"`
	DocumentDate     string  `xml:"Document_Date,omitempty" json:"documentDate,omitempty"`
	DocumentType     string  `xml:"Document_Type,omitempty" json:"documentType,omitempty"`
	AccountType      string  `xml:"Account_Type,omitempty" json:"accountType,omitempty"`
	AccountNo        string  `xml:"Account_No,omitempty" json:"accountNo,omitempty"`
	Amount           float64 `xml:"Amount,omitempty" json:"amount,omitempty"`
	AppliesToDocType string  `xml:"Applies_to_Doc_Type,omitempty" json:"appliesToDocType,omitempty"`
	AppliesToDocNo   string  `xml:"Applies_to_Doc_No,omitempty" json:"appliesToDocNo,omitempty"`
}

type BackToCDSLedgerEntriesResponse struct {
	VendorNo   string `json:"vendor_no,omitempty"`
	DocumentNo int    `json:"document_no,omitempty"`
}
