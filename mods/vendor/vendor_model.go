package vendor

// type AddVendorModel struct {
// 	WSVendor struct {
// 		// Key                             string `xml:"Key" json:"key"`
// 		// No                              string `xml:"No" json:"no"`
// 		Name              string `xml:"Name" json:"name"`
// 		Address           string `xml:"Address" json:"address"`
// 		Address2          string `xml:"Address_2" json:"address_2"`
// 		PostCode          string `xml:"Post_Code" json:"post_code"`
// 		City              string `xml:"City" json:"city"`
// 		County            string `xml:"County" json:"county"`
// 		CountryRegionCode string `xml:"Country_Region_Code" json:"country_region_code"`
// 		PhoneNo           string `xml:"Phone_No" json:"phone_no"`
// 		PrimaryContactNo  string `xml:"Primary_Contact_No" json:"primary_contact_no"`
// 		Contact           string `xml:"Contact" json:"contact"`
// 		// SearchName                      string `xml:"Search_Name" json:"search_name"`
// 		// BalanceLCY                      string `xml:"Balance_LCY" json:"balance_lcy"`
// 		// PostDatedChecksLCY              string `xml:"Post_Dated_Checks_LCY" json:"post_dated_checks_lcy"`
// 		// BalanceLCYABSPostDatedChecksLCY string `xml:"_Balance_LCY_ABS_Post_Dated_Checks_LCY" json:"balance_lcy_abs_post_dated_checks_lcy"`
// 		// PurchaserCode                   string `xml:"Purchaser_Code" json:"purchaser_code"`
// 		// ResponsibilityCenter            string `xml:"Responsibility_Center" json:"responsibility_center"`
// 		// Blocked                         string `xml:"Blocked" json:"blocked"`
// 		// BankAccountModifiedBy           string `xml:"Bank_Account_Modified_By" json:"bank_account_modified_by"`
// 		// LastDateModified                string `xml:"Last_Date_Modified" json:"last_date_modified"`
// 		WeighbridgeSupplierID string `xml:"Weighbridge_Supplier_ID" json:"weighbridge_supplier_id"`
// 		// FaxNo                           string `xml:"Fax_No" json:"fax_no"`
// 		// EMail                           string `xml:"E_Mail" json:"e_mail"`
// 		// HomePage                        string `xml:"Home_Page" json:"home_page"`
// 		// ICPartnerCode                   string `xml:"IC_Partner_Code" json:"ic_partner_code"`
// 		// PayToVendorNo                   string `xml:"Pay_to_Vendor_No" json:"pay_to_vendor_no"`
// 		// GenBusPostingGroup              string `xml:"Gen_Bus_Posting_Group" json:"gen_bus_posting_group"`
// 		// VATBusPostingGroup              string `xml:"VAT_Bus_Posting_Group" json:"vat_bus_posting_group"`
// 		// WHTBusinessPostingGroup         string `xml:"WHT_Business_Posting_Group" json:"wht_business_posting_group"`
// 		// VendorPostingGroup              string `xml:"Vendor_Posting_Group" json:"vendor_posting_group"`
// 		// InvoiceDiscCode                 string `xml:"Invoice_Disc_Code" json:"invoice_disc_code"`
// 		// PricesIncludingVAT              string `xml:"Prices_Including_VAT" json:"prices_including_vat"`
// 		// PrepaymentPercent               string `xml:"Prepayment_Percent" json:"prepayment_percent"`
// 		// ApplicationMethod               string `xml:"Application_Method" json:"application_method"`
// 		// PartnerType                     string `xml:"Partner_Type" json:"partner_type"`
// 		// PaymentTermsCode                string `xml:"Payment_Terms_Code" json:"payment_terms_code"`
// 		// PaymentMethodCode               string `xml:"Payment_Method_Code" json:"payment_method_code"`
// 		// Priority                        string `xml:"Priority" json:"priority"`
// 		// CashFlowPaymentTermsCode        string `xml:"Cash_Flow_Payment_Terms_Code" json:"cash_flow_payment_terms_code"`
// 		// LodgementReference              string `xml:"Lodgement_Reference" json:"lodgement_reference"`
// 		// OurAccountNo                    string `xml:"Our_Account_No" json:"our_account_no"`
// 		// BlockPaymentTolerance           string `xml:"Block_Payment_Tolerance" json:"block_payment_tolerance"`
// 		// EFTPayment                      string `xml:"EFT_Payment" json:"eft_payment"`
// 		// EFTBankAccountNo                string `xml:"EFT_Bank_Account_No" json:"eft_bank_account_no"`
// 		// CreditorNo                      string `xml:"Creditor_No" json:"creditor_no"`
// 		// PreferredBankAccount            string `xml:"Preferred_Bank_Account" json:"preferred_bank_account"`
// 		// LocationCode                    string `xml:"Location_Code" json:"location_code"`
// 		// ShipmentMethodCode              string `xml:"Shipment_Method_Code" json:"shipment_method_code"`
// 		// LeadTimeCalculation             string `xml:"Lead_Time_Calculation" json:"lead_time_calculation"`
// 		// BaseCalendarCode                string `xml:"Base_Calendar_Code" json:"base_calendar_code"`
// 		// CustomizedCalendar              string `xml:"Customized_Calendar" json:"customized_calendar"`
// 		// CurrencyCode                    string `xml:"Currency_Code" json:"currency_code"`
// 		// LanguageCode                    string `xml:"Language_Code" json:"language_code"`
// 		// VATRegistrationNo               string `xml:"VAT_Registration_No" json:"vat_registration_no"`
// 		// WHTRegistrationID               string `xml:"WHT_Registration_ID" json:"wht_registration_id"`
// 		// IDNo                            string `xml:"ID_No" json:"id_no"`
// 		// IRDNo                           string `xml:"IRD_No" json:"ird_no"`
// 		// ABN                             string `xml:"ABN" json:"abn"`
// 		// ABNDivisionPartNo               string `xml:"ABN_Division_Part_No" json:"abn_division_part_no"`
// 		// Registered                      string `xml:"Registered" json:"registered"`
// 		// ForeignVend                     string `xml:"Foreign_Vend" json:"foreign_vend"`
// 		// SendToConcur                    string `xml:"Send_To_Concur" json:"send_to_concur"`
// 		// ConcurInvoiceLastUpdated        string `xml:"Concur_Invoice_Last_Updated" json:"concur_invoice_last_updated"`
// 	} `xml:"WSVendor" json:"ws_vendor"`
// }

type WSVendor struct {
	// Key                             string `xml:"Key" json:"key"`
	// No                              string `xml:"No" json:"no"`
	Name     string `xml:"Name" json:"name"`
	Address  string `xml:"Address" json:"address"`
	Address2 string `xml:"Address_2" json:"address_2"`
	PostCode string `xml:"Post_Code" json:"post_code"`
	City     string `xml:"City" json:"city"`
	County   string `xml:"County" json:"county"`
	// CountryRegionCode               string `xml:"Country_Region_Code" json:"country_region_code"`
	// PhoneNo                         string `xml:"Phone_No" json:"phone_no"`
	// PrimaryContactNo                string `xml:"Primary_Contact_No" json:"primary_contact_no"`
	// Contact                         string `xml:"Contact" json:"contact"`
	// SearchName                      string `xml:"Search_Name" json:"search_name"`
	// BalanceLCY                      string `xml:"Balance_LCY" json:"balance_lcy"`
	// PostDatedChecksLCY              string `xml:"Post_Dated_Checks_LCY" json:"post_dated_checks_lcy"`
	// BalanceLCYABSPostDatedChecksLCY string `xml:"_Balance_LCY_ABS_Post_Dated_Checks_LCY" json:"balance_lcy_abs_post_dated_checks_lcy"`
	// PurchaserCode                   string `xml:"Purchaser_Code" json:"purchaser_code"`
	// ResponsibilityCenter            string `xml:"Responsibility_Center" json:"responsibility_center"`
	// Blocked                         string `xml:"Blocked" json:"blocked"`
	// BankAccountModifiedBy           string `xml:"Bank_Account_Modified_By" json:"bank_account_modified_by"`
	// LastDateModified                string `xml:"Last_Date_Modified" json:"last_date_modified"`
	WeighbridgeSupplierID string `xml:"Weighbridge_Supplier_ID" json:"weighbridge_supplier_id"`
	// FaxNo                           string `xml:"Fax_No" json:"fax_no"`
	// EMail                           string `xml:"E_Mail" json:"e_mail"`
	// HomePage                        string `xml:"Home_Page" json:"home_page"`
	// ICPartnerCode                   string `xml:"IC_Partner_Code" json:"ic_partner_code"`
	// PayToVendorNo                   string `xml:"Pay_to_Vendor_No" json:"pay_to_vendor_no"`
	// GenBusPostingGroup              string `xml:"Gen_Bus_Posting_Group" json:"gen_bus_posting_group"`
	// VATBusPostingGroup              string `xml:"VAT_Bus_Posting_Group" json:"vat_bus_posting_group"`
	// WHTBusinessPostingGroup         string `xml:"WHT_Business_Posting_Group" json:"wht_business_posting_group"`
	// VendorPostingGroup              string `xml:"Vendor_Posting_Group" json:"vendor_posting_group"`
	// InvoiceDiscCode                 string `xml:"Invoice_Disc_Code" json:"invoice_disc_code"`
	// PricesIncludingVAT              string `xml:"Prices_Including_VAT" json:"prices_including_vat"`
	// PrepaymentPercent               string `xml:"Prepayment_Percent" json:"prepayment_percent"`
	// ApplicationMethod               string `xml:"Application_Method" json:"application_method"`
	// PartnerType                     string `xml:"Partner_Type" json:"partner_type"`
	// PaymentTermsCode                string `xml:"Payment_Terms_Code" json:"payment_terms_code"`
	// PaymentMethodCode               string `xml:"Payment_Method_Code" json:"payment_method_code"`
	// Priority                        string `xml:"Priority" json:"priority"`
	// CashFlowPaymentTermsCode        string `xml:"Cash_Flow_Payment_Terms_Code" json:"cash_flow_payment_terms_code"`
	// LodgementReference              string `xml:"Lodgement_Reference" json:"lodgement_reference"`
	// OurAccountNo                    string `xml:"Our_Account_No" json:"our_account_no"`
	// BlockPaymentTolerance           string `xml:"Block_Payment_Tolerance" json:"block_payment_tolerance"`
	// EFTPayment                      string `xml:"EFT_Payment" json:"eft_payment"`
	// EFTBankAccountNo                string `xml:"EFT_Bank_Account_No" json:"eft_bank_account_no"`
	// CreditorNo                      string `xml:"Creditor_No" json:"creditor_no"`
	// PreferredBankAccount            string `xml:"Preferred_Bank_Account" json:"preferred_bank_account"`
	// LocationCode                    string `xml:"Location_Code" json:"location_code"`
	// ShipmentMethodCode              string `xml:"Shipment_Method_Code" json:"shipment_method_code"`
	// LeadTimeCalculation             string `xml:"Lead_Time_Calculation" json:"lead_time_calculation"`
	// BaseCalendarCode                string `xml:"Base_Calendar_Code" json:"base_calendar_code"`
	// CustomizedCalendar              string `xml:"Customized_Calendar" json:"customized_calendar"`
	// CurrencyCode                    string `xml:"Currency_Code" json:"currency_code"`
	// LanguageCode                    string `xml:"Language_Code" json:"language_code"`
	// VATRegistrationNo               string `xml:"VAT_Registration_No" json:"vat_registration_no"`
	// WHTRegistrationID               string `xml:"WHT_Registration_ID" json:"wht_registration_id"`
	// IDNo                            string `xml:"ID_No" json:"id_no"`
	// IRDNo                           string `xml:"IRD_No" json:"ird_no"`
	// ABN                             string `xml:"ABN" json:"abn"`
	// ABNDivisionPartNo               string `xml:"ABN_Division_Part_No" json:"abn_division_part_no"`
	// Registered                      string `xml:"Registered" json:"registered"`
	// ForeignVend                     string `xml:"Foreign_Vend" json:"foreign_vend"`
	// SendToConcur                    string `xml:"Send_To_Concur" json:"send_to_concur"`
	// ConcurInvoiceLastUpdated        string `xml:"Concur_Invoice_Last_Updated" json:"concur_invoice_last_updated"`
}

//`xml:"WSVendor" json:"ws_vendor"`
