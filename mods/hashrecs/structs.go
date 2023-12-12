package hashrecs

type HashRec struct {
	Hash       string `json:"hash,omitempty"`
	NavID      string `json:"nav_id,omitempty"`
	InvoiceNo  string `json:"invoice_no,omitempty"`
	DocumentNo string `json:"document_no,omitempty"`
	RefundId   int    `json:"refund_id,omitempty"`
	PaymentId  int    `json:"payment_id,omitempty"`
}

type HashRecs struct {
	FilePath string
	Name     string
	Recs     map[string]HashRec
}
