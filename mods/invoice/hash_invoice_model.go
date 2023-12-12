package invoice

//Here key string is RefundNo -> Vendor_Invoice_No
type HashInvoiceModel map[string]HashInvoiceEntry

type HashInvoiceEntry struct {
	Hash       string `json:"hash"`
	NavID      string `json:"nav_id"`
	InvoiceNo  string `json:"invoice_no"`
	DocumentNo string `json:"document_no"`
	RefundId   int    `json:"refund_id"`
}
