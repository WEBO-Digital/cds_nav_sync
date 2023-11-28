package invoice

// WSPurchaseInvoicePage represents the structure of WSPurchaseInvoicePage element in the XML
type WSPurchaseInvoicePage2 struct {
	Key               string `xml:"Key" json:"key"`
	No                string `xml:"No" json:"no"`
	BuyFromVendorNo   string `xml:"Buy_from_Vendor_No" json:"buy_from_vendor_no"`
	BuyFromContactNo  string `xml:"Buy_from_Contact_No" json:"buy_from_contact_no"`
	BuyFromVendorName string `xml:"Buy_from_Vendor_Name" json:"buy_from_vendor_name"`
	// Add other fields based on your XML structure
	PurchLines PurchLines `xml:"PurchLines" json:"purch_lines"`
}

// PurchLines represents the structure of PurchLines element in the XML
type PurchLines struct {
	PurchInvoiceLine PurchInvoiceLine `xml:"Purch_Invoice_Line" json:"purch_invoice_line"`
}

// PurchInvoiceLine represents the structure of Purch_Invoice_Line element in the XML
type PurchInvoiceLine struct {
	Key  string `xml:"Key" json:"key"`
	Type string `xml:"Type" json:"type"`
	// Add other fields based on your XML structure
}

// Envelope represents the structure of Soap:Envelope element in the XML
type PostInvoiceEnvelope struct {
	Body Body `xml:"Body" json:"body"`
}

// Body represents the structure of Soap:Body element in the XML
type Body struct {
	CreateResult CreateResult `xml:"Create_Result" json:"create_result"`
}

// CreateResult represents the structure of Create_Result element in the XML
type CreateResult struct {
	WSPurchaseInvoicePage WSPurchaseInvoicePage2 `xml:"WSPurchaseInvoicePage" json:"ws_purchase_invoice_page"`
}
