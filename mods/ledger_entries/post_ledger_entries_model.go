package ledgerentries

type PostLedgerEntriesEnvelope struct {
	Body Body `xml:"http://schemas.xmlsoap.org/soap/envelope/ Body" json:"body"`
}

type Body struct {
	CreateResult CreateResult `xml:"urn:microsoft-dynamics-schemas/page/vendorpayment Create_Result" json:"create_result"`
}

type CreateResult struct {
	VendorPayment LEVendorPayment `xml:"VendorPayment" json:"vendor_payment"`
}

type LEVendorPayment struct {
	Key                     string  `xml:"Key" json:"key"`
	PostingDate             string  `xml:"Posting_Date" json:"posting_date"`
	DocumentDate            string  `xml:"Document_Date" json:"document_date"`
	DocumentType            string  `xml:"Document_Type" json:"document_type"`
	DocumentNo              int     `xml:"Document_No" json:"document_no"`
	IncomingDocumentEntryNo int     `xml:"Incoming_Document_Entry_No" json:"incoming_document_entry_no"`
	AppliesToExtDocNo       string  `xml:"Applies_to_Ext_Doc_No" json:"applies_to_ext_doc_no"`
	AccountType             string  `xml:"Account_Type" json:"account_type"`
	AccountNo               string  `xml:"Account_No" json:"account_no"`
	Description             string  `xml:"Description" json:"description"`
	GenPostingType          string  `xml:"Gen_Posting_Type" json:"gen_posting_type"`
	WHTPayment              bool    `xml:"WHT_Payment" json:"wht_payment"`
	SkipWHT                 bool    `xml:"Skip_WHT" json:"skip_wht"`
	Amount                  float32 `xml:"Amount" json:"amount"`
	AmountLCY               float32 `xml:"Amount_LCY" json:"amount_lcy"`
	DebitAmount             float32 `xml:"Debit_Amount" json:"debit_amount"`
	CreditAmount            float32 `xml:"Credit_Amount" json:"credit_amount"`
	VATAmount               float32 `xml:"VAT_Amount" json:"vat_amount"`
	VATDifference           float32 `xml:"VAT_Difference" json:"vat_difference"`
	VendorExchangeRateACY   float32 `xml:"Vendor_Exchange_Rate_ACY" json:"vendor_exchange_rate_acy"`
	BalVATAmount            float32 `xml:"Bal_VAT_Amount" json:"bal_vat_amount"`
	BalVATDifference        float32 `xml:"Bal_VAT_Difference" json:"bal_vat_difference"`
	BalAccountType          string  `xml:"Bal_Account_Type" json:"bal_account_type"`
	BalAccountNo            string  `xml:"Bal_Account_No" json:"bal_account_no"`
	BalGenPostingType       string  `xml:"Bal_Gen_Posting_Type" json:"bal_gen_posting_type"`
	ShortcutDimension1Code  string  `xml:"Shortcut_Dimension_1_Code" json:"shortcut_dimension_1_code"`
	ShortcutDimension2Code  string  `xml:"Shortcut_Dimension_2_Code" json:"shortcut_dimension_2_code"`
	AppliedYesNo            bool    `xml:"Applied_Yes_No" json:"applied_yes_no"`
	AppliesToDocType        string  `xml:"Applies_to_Doc_Type" json:"applies_to_doc_type"`
	AppliesToDocNo          string  `xml:"Applies_to_Doc_No" json:"applies_to_doc_no"`
	GetAppliesToDocDueDate  string  `xml:"GetAppliesToDocDueDate" json:"get_applies_to_doc_due_date"`
	BankPaymentType         string  `xml:"Bank_Payment_Type" json:"bank_payment_type"`
	CheckPrinted            bool    `xml:"Check_Printed" json:"check_printed"`
	OverdueWarningText      string  `xml:"OverdueWarningText" json:"overdue_warning_text"`
	AccName                 string  `xml:"AccName" json:"acc_name"`
	BalAccName              string  `xml:"BalAccName" json:"bal_acc_name"`
	Balance                 float32 `xml:"Balance" json:"balance"`
	TotalBalance            float32 `xml:"TotalBalance" json:"total_balance"`
}
