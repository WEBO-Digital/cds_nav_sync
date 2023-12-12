package invoice

type WSPurchaseInvoicePage struct {
	RefundId                 int           `json:"refund_id,omitempty"`
	BuyFromVendorNo          *string       `xml:"Buy_from_Vendor_No,omitempty" json:"buy_from_vendor_no,omitempty"`
	BuyFromContactNo         string        `xml:"Buy_from_Contact_No,omitempty" json:"buy_from_contact_no,omitempty"`
	BuyFromVendorName        string        `xml:"Buy_from_Vendor_Name,omitempty" json:"buy_from_vendor_name,omitempty"`
	BuyFromCity              string        `xml:"Buy_from_City,omitempty" json:"buy_from_city,omitempty"`
	BuyFromCountryRegionCode string        `xml:"Buy_from_Country_Region_Code,omitempty" json:"buy_from_country_region_code,omitempty"`
	PostingDate              string        `xml:"Posting_Date,omitempty" json:"posting_date,omitempty"`
	DocumentDate             string        `xml:"Document_Date,omitempty" json:"document_date,omitempty"`
	VendorInvoiceNo          string        `xml:"Vendor_Invoice_No,omitempty" json:"vendor_invoice_no,omitempty"`
	PayToVendorNo            string        `xml:"Pay_to_Vendor_No,omitempty" json:"pay_to_vendor_no,omitempty"`
	PayToName                string        `xml:"Pay_to_Name,omitempty" json:"pay_to_name,omitempty"`
	PayToCity                string        `xml:"Pay_to_City,omitempty" json:"pay_to_city,omitempty"`
	PayToCountryRegionCode   string        `xml:"Pay_to_Country_Region_Code,omitempty" json:"pay_to_country_region_code,omitempty"`
	PaymentTermsCode         string        `xml:"Payment_Terms_Code,omitempty" json:"payment_terms_code,omitempty"`
	DueDate                  string        `xml:"Due_Date,omitempty" json:"due_date,omitempty"`
	VATBusPostingGroup       string        `xml:"VAT_Bus_Posting_Group,omitempty" json:"vat_bus_posting_group,omitempty"`
	ShipToName               string        `xml:"Ship_to_Name,omitempty" json:"ship_to_name,omitempty"`
	ShipToAddress            string        `xml:"Ship_to_Address,omitempty" json:"ship_to_address,omitempty"`
	ShipToAddress2           string        `xml:"Ship_to_Address_2,omitempty" json:"ship_to_address_2,omitempty"`
	ShipToPostCode           string        `xml:"Ship_to_Post_Code,omitempty" json:"ship_to_post_code,omitempty"`
	ShipToCity               string        `xml:"Ship_to_City,omitempty" json:"ship_to_city,omitempty"`
	ShipToCountryRegionCode  string        `xml:"Ship_to_Country_Region_Code,omitempty" json:"ship_to_country_region_code,omitempty"`
	LocationCode             string        `xml:"Location_Code,omitempty" json:"location_code,omitempty"`
	ExpectedReceiptDate      string        `xml:"Expected_Receipt_Date,omitempty" json:"expected_receipt_date,omitempty"`
	PurchLines               []PurchLines2 `xml:"PurchLines,omitempty" json:"purch_lines,omitempty"`
}

type PurchLines2 struct {
	PurchInvoiceLine PurchInvoiceLine `xml:"Purch_Invoice_Line,omitempty" json:"purch_invoice_line,omitempty"`
}

type PurchInvoiceLine struct {
	Type               string `xml:"Type,omitempty" json:"type,omitempty"`
	No                 string `xml:"No,omitempty" json:"no,omitempty"`
	Quantity           string `xml:"Quantity,omitempty" json:"quantity,omitempty"`
	ShortcutDimension1 string `xml:"Shortcut_Dimension_1_Code,omitempty" json:"shortcut_dimension_1_code,omitempty"`
	ShortcutDimension2 string `xml:"Shortcut_Dimension_2_Code,omitempty" json:"shortcut_dimension_2_code,omitempty"`
	UnitPriceLCY       string `xml:"Unit_Price_LCY,omitempty" json:"unit_price_lcy,omitempty"`
	LocationCode       string `xml:"Location_Code,omitempty" json:"location_code,omitempty"`
}

type BackToCDSInvoiceResponse struct {
	RefundId          string `json:"refund_id,omitempty"`
	VendorNo          string `json:"vendor_no,omitempty"`
	PurchaseInvoiceNo string `json:"purchase_invoice_no,omitempty"`
	DocumentNo        string `json:"document_id,omitempty"`
}

// type WSPurchaseInvoicePage struct {
// 	BuyFromVendorNo string `xml:"Buy_from_Vendor_No" json:"buy_from_vendor_no"`
// 	// BuyFromContactNo         string `xml:"Buy_from_Contact_No" json:"buy_from_contact_no"`
// 	BuyFromVendorName string `xml:"Buy_from_Vendor_Name" json:"buy_from_vendor_name"`
// 	// BuyFromAddress           string `xml:"Buy_from_Address" json:"buy_from_address"`
// 	// BuyFromAddress2          string `xml:"Buy_from_Address_2" json:"buy_from_address_2"`
// 	// BuyFromPostCode          string `xml:"Buy_from_Post_Code" json:"buy_from_post_code"`
// 	// BuyFromCity              string `xml:"Buy_from_City" json:"buy_from_city"`
// 	// BuyFromCounty            string `xml:"Buy_from_County" json:"buy_from_county"`
// 	// BuyFromCountryRegionCode string `xml:"Buy_from_Country_Region_Code" json:"buy_from_country_region_code"`
// 	// BuyFromContact           string `xml:"Buy_from_Contact" json:"buy_from_contact"`
// 	// PostingDate              string `xml:"Posting_Date" json:"posting_date"`
// 	// DocumentDate             string `xml:"Document_Date" json:"document_date"`
// 	// IncomingDocumentEntryNo  string `xml:"Incoming_Document_Entry_No" json:"incoming_document_entry_no"`
// 	VendorInvoiceNo string `xml:"Vendor_Invoice_No" json:"vendor_invoice_no"`
// 	// OrderAddressCode         string `xml:"Order_Address_Code" json:"order_address_code"`
// 	PurchaserCode string `xml:"Purchaser_Code" json:"purchaser_code"`
// 	// ConcurInvoiceApprovalMgr string `xml:"Concur_Invoice_Approval_Mgr" json:"concur_invoice_approval_mgr"`
// 	// CampaignNo               string `xml:"Campaign_No" json:"campaign_no"`
// 	// ResponsibilityCenter     string `xml:"Responsibility_Center" json:"responsibility_center"`
// 	// AssignedUserID           string `xml:"Assigned_User_ID" json:"assigned_user_id"`
// 	// JobQueueStatus           string `xml:"Job_Queue_Status" json:"job_queue_status"`
// 	// Status                   string `xml:"Status" json:"status"`
// 	// XaanaTransaction         string `xml:"Xaana_Transaction" json:"xaana_transaction"`
// 	// PayToVendorNo            string `xml:"Pay_to_Vendor_No" json:"pay_to_vendor_no"`
// 	// PayToContactNo           string `xml:"Pay_to_Contact_No" json:"pay_to_contact_no"`
// 	// PayToName                string `xml:"Pay_to_Name" json:"pay_to_name"`
// 	// PayToAddress             string `xml:"Pay_to_Address" json:"pay_to_address"`
// 	// PayToAddress2            string `xml:"Pay_to_Address_2" json:"pay_to_address_2"`
// 	// PayToPostCode            string `xml:"Pay_to_Post_Code" json:"pay_to_post_code"`
// 	// PayToCity                string `xml:"Pay_to_City" json:"pay_to_city"`
// 	// PayToCounty              string `xml:"Pay_to_County" json:"pay_to_county"`
// 	// PayToCountryRegionCode   string `xml:"Pay_to_Country_Region_Code" json:"pay_to_country_region_code"`
// 	// PayToContact             string `xml:"Pay_to_Contact" json:"pay_to_contact"`
// 	// VendorExchangeRateACY    string `xml:"Vendor_Exchange_Rate_ACY" json:"vendor_exchange_rate_acy"`
// 	// ShortcutDimension1Code string `xml:"Shortcut_Dimension_1_Code,omitempty" json:"shortcut_dimension_1_code,omitempty"`
// 	// ShortcutDimension2Code string `xml:"Shortcut_Dimension_2_Code,omitempty" json:"shortcut_dimension_2_code,omitempty"`
// 	// PaymentTermsCode         string `xml:"Payment_Terms_Code" json:"payment_terms_code"`
// 	// DueDate                  string `xml:"Due_Date" json:"due_date"`
// 	// PaymentDiscountPercent   string `xml:"Payment_Discount_Percent" json:"payment_discount_percent"`
// 	// PmtDiscountDate          string `xml:"Pmt_Discount_Date" json:"pmt_discount_date"`
// 	// PaymentMethodCode        string `xml:"Payment_Method_Code" json:"payment_method_code"`
// 	// PaymentReference         string `xml:"Payment_Reference" json:"payment_reference"`
// 	// CreditorNo               string `xml:"Creditor_No" json:"creditor_no"`
// 	// OnHold                   string `xml:"On_Hold" json:"on_hold"`
// 	// PricesIncludingVAT       string `xml:"Prices_Including_VAT" json:"prices_including_vat"`
// 	// VATBusPostingGroup       string `xml:"VAT_Bus_Posting_Group" json:"vat_bus_posting_group"`
// 	// InvoiceReceivedDate      string `xml:"Invoice_Received_Date" json:"invoice_received_date"`
// ShipToName               string `xml:"Ship_to_Name" json:"ship_to_name"`
// ShipToAddress            string `xml:"Ship_to_Address" json:"ship_to_address"`
// ShipToAddress2           string `xml:"Ship_to_Address_2" json:"ship_to_address_2"`
// ShipToPostCode           string `xml:"Ship_to_Post_Code" json:"ship_to_post_code"`
// ShipToCity               string `xml:"Ship_to_City" json:"ship_to_city"`
// ShipToCounty             string `xml:"Ship_to_County" json:"ship_to_county"`
// ShipToCountryRegionCode  string `xml:"Ship_to_Country_Region_Code" json:"ship_to_country_region_code"`
// 	// ShipToContact            string `xml:"Ship_to_Contact" json:"ship_to_contact"`
// 	// LocationCode             string `xml:"Location_Code" json:"location_code"`
// 	// ShipmentMethodCode       string `xml:"Shipment_Method_Code" json:"shipment_method_code"`
// 	// ExpectedReceiptDate      string `xml:"Expected_Receipt_Date" json:"expected_receipt_date"`
// 	// CurrencyCode             string `xml:"Currency_Code" json:"currency_code"`
// 	// TransactionType          string `xml:"Transaction_Type" json:"transaction_type"`
// 	// TransactionSpecification string `xml:"Transaction_Specification" json:"transaction_specification"`
// 	// TransportMethod          string `xml:"Transport_Method" json:"transport_method"`
// 	// EntryPoint               string `xml:"Entry_Point" json:"entry_point"`
// 	// Area                     string `xml:"Area" json:"area"`
// 	// AppliesToDocType         string `xml:"Applies_to_Doc_Type" json:"applies_to_doc_type"`
// 	// AppliesToDocNo           string `xml:"Applies_to_Doc_No" json:"applies_to_doc_no"`
// 	PurchLines []struct {
// 		PurchInvoiceLine struct {
// 			Type string `xml:"Type" json:"type"`
// 			No   string `xml:"No" json:"no"`
// 			// CrossReferenceNo            string `xml:"Cross_Reference_No" json:"cross_reference_no"`
// 			// GenProdPostingGroup         string `xml:"Gen_Prod_Posting_Group" json:"gen_prod_posting_group"`
// 			// FAPostingType               string `xml:"FA_Posting_Type" json:"fa_posting_type"`
// 			// GenBusPostingGroup          string `xml:"Gen_Bus_Posting_Group" json:"gen_bus_posting_group"`
// 			// MaintenanceCode             string `xml:"Maintenance_Code" json:"maintenance_code"`
// 			// ICPartnerCode               string `xml:"IC_Partner_Code" json:"ic_partner_code"`
// 			// ICPartnerRefType            string `xml:"IC_Partner_Ref_Type" json:"ic_partner_ref_type"`
// 			// ICPartnerReference          string `xml:"IC_Partner_Reference" json:"ic_partner_reference"`
// 			// VariantCode                 string `xml:"Variant_Code" json:"variant_code"`
// 			// Nonstock                    string `xml:"Nonstock" json:"nonstock"`
// 			// VATProdPostingGroup         string `xml:"VAT_Prod_Posting_Group" json:"vat_prod_posting_group"`
// 			// WHTBusinessPostingGroup     string `xml:"WHT_Business_Posting_Group" json:"wht_business_posting_group"`
// 			// WHTProductPostingGroup      string `xml:"WHT_Product_Posting_Group" json:"wht_product_posting_group"`
// 			// Description                 string `xml:"Description" json:"description"`
// 			// VendorShipmentNo            string `xml:"Vendor_Shipment_No" json:"vendor_shipment_no"`
// 			// ApportionMethod             string `xml:"Apportion_Method" json:"apportion_method"`
// 			// ReturnReasonCode            string `xml:"Return_Reason_Code" json:"return_reason_code"`
// 			LocationCode string `xml:"Location_Code" json:"location_code"`
// 			// BinCode                     string `xml:"Bin_Code" json:"bin_code"`
// 			Quantity string `xml:"Quantity" json:"quantity"`
// 			// UnitOfMeasureCode           string `xml:"Unit_of_Measure_Code" json:"unit_of_measure_code"`
// 			// UnitOfMeasure               string `xml:"Unit_of_Measure" json:"unit_of_measure"`
// 			// DirectUnitCost              string `xml:"Direct_Unit_Cost" json:"direct_unit_cost"`
// 			// IndirectCostPercent         string `xml:"Indirect_Cost_Percent" json:"indirect_cost_percent"`
// 			// UnitCostLCY                 string `xml:"Unit_Cost_LCY" json:"unit_cost_lcy"`
// 			UnitPriceLCY string `xml:"Unit_Price_LCY" json:"unit_price_lcy"`
// 			// LineAmount                  string `xml:"Line_Amount" json:"line_amount"`
// 			// WHTAbsorbBase               string `xml:"WHT_Absorb_Base" json:"wht_absorb_base"`
// 			// LineDiscountPercent         string `xml:"Line_Discount_Percent" json:"line_discount_percent"`
// 			// LineDiscountAmount          string `xml:"Line_Discount_Amount" json:"line_discount_amount"`
// 			// AllowInvoiceDisc            string `xml:"Allow_Invoice_Disc" json:"allow_invoice_disc"`
// 			// InvDiscountAmount           string `xml:"Inv_Discount_Amount" json:"inv_discount_amount"`
// 			// AllowItemChargeAssignment   string `xml:"Allow_Item_Charge_Assignment" json:"allow_item_charge_assignment"`
// 			// QtyToAssign                 string `xml:"Qty_to_Assign" json:"qty_to_assign"`
// 			// QtyAssigned                 string `xml:"Qty_Assigned" json:"qty_assigned"`
// 			// JobNo                       string `xml:"Job_No" json:"job_no"`
// 			// JobTaskNo                   string `xml:"Job_Task_No" json:"job_task_no"`
// 			// JobLineType                 string `xml:"Job_Line_Type" json:"job_line_type"`
// 			// JobUnitPrice                string `xml:"Job_Unit_Price" json:"job_unit_price"`
// 			// JobLineAmount               string `xml:"Job_Line_Amount" json:"job_line_amount"`
// 			// JobLineDiscountAmount       string `xml:"Job_Line_Discount_Amount" json:"job_line_discount_amount"`
// 			// JobLineDiscountPercent      string `xml:"Job_Line_Discount_Percent" json:"job_line_discount_percent"`
// 			// JobTotalPrice               string `xml:"Job_Total_Price" json:"job_total_price"`
// 			// JobUnitPriceLCY             string `xml:"Job_Unit_Price_LCY" json:"job_unit_price_lcy"`
// 			// JobTotalPriceLCY            string `xml:"Job_Total_Price_LCY" json:"job_total_price_lcy"`
// 			// JobLineAmountLCY            string `xml:"Job_Line_Amount_LCY" json:"job_line_amount_lcy"`
// 			// JobLineDiscAmountLCY        string `xml:"Job_Line_Disc_Amount_LCY" json:"job_line_disc_amount_lcy"`
// 			// ProdOrderNo                 string `xml:"Prod_Order_No" json:"prod_order_no"`
// 			// BlanketOrderNo              string `xml:"Blanket_Order_No" json:"blanket_order_no"`
// 			// BlanketOrderLineNo          string `xml:"Blanket_Order_Line_No" json:"blanket_order_line_no"`
// 			// InsuranceNo                 string `xml:"Insurance_No" json:"insurance_no"`
// 			// BudgetedFANo                string `xml:"Budgeted_FA_No" json:"budgeted_fa_no"`
// 			// DepreciationBookCode        string `xml:"Depreciation_Book_Code" json:"depreciation_book_code"`
// 			// DeprUntilFAPostingDate      string `xml:"Depr_until_FA_Posting_Date" json:"depr_until_fa_posting_date"`
// 			// DeprAcquisitionCost         string `xml:"Depr_Acquisition_Cost" json:"depr_acquisition_cost"`
// 			// DuplicateInDepreciationBook string `xml:"Duplicate_in_Depreciation_Book" json:"duplicate_in_depreciation_book"`
// 			// UseDuplicationList          string `xml:"Use_Duplication_List" json:"use_duplication_list"`
// 			// ApplToItemEntry             string `xml:"Appl_to_Item_Entry" json:"appl_to_item_entry"`
// 			// ShortcutDimension1Code      string `xml:"Shortcut_Dimension_1_Code" json:"shortcut_dimension_1_code"`
// 			// ShortcutDimension2Code      string `xml:"Shortcut_Dimension_2_Code" json:"shortcut_dimension_2_code"`
// 			// ShortcutDimCodeX005B3X005D  string `xml:"ShortcutDimCode_x005B_3_x005D_" json:"shortcut_dim_code_x005b_3_x005d_"`
// 			// ShortcutDimCodeX005B4X005D  string `xml:"ShortcutDimCode_x005B_4_x005D_" json:"shortcut_dim_code_x005b_4_x005d_"`
// 			// ShortcutDimCodeX005B5X005D  string `xml:"ShortcutDimCode_x005B_5_x005D_" json:"shortcut_dim_code_x005b_5_x005d_"`
// 			// ShortcutDimCodeX005B6X005D  string `xml:"ShortcutDimCode_x005B_6_x005D_" json:"shortcut_dim_code_x005b_6_x005d_"`
// 			// ShortcutDimCodeX005B7X005D  string `xml:"ShortcutDimCode_x005B_7_x005D_" json:"shortcut_dim_code_x005b_7_x005d_"`
// 			// ShortcutDimCodeX005B8X005D  string `xml:"ShortcutDimCode_x005B_8_x005D_" json:"shortcut_dim_code_x005b_8_x005d_"`
// 			// TicketLineNo                string `xml:"Ticket_Line_No" json:"ticket_line_no"`
// 		} `xml:"Purch_Invoice_Line" json:"purch_invoice_line"`
// 	} `xml:"PurchLines" json:"purch_lines"`
// }
