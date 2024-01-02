package vendor

import (
	"encoding/xml"
)

type CreateResultVendor struct {
	XMLName xml.Name   `xml:"Envelope,omitempty" json:"envelope,omitempty"`
	Body    BodyVendor `xml:"Body,omitempty" json:"body,omitempty"`
}

type BodyVendor struct {
	CreateResult CreateResultDetailVendor `xml:"Create_Result,omitempty" json:"create_result,omitempty"`
}

type CreateResultDetailVendor struct {
	WSVendor WSVendor `xml:"WSVendor,omitempty" json:"ws_vendor,omitempty"`
}

type WSVendor struct {
	Key               *string `xml:"Key,omitempty" json:"key,omitempty"`
	No                *string `xml:"No,omitempty" json:"no,omitempty"`
	Name              string  `xml:"Name,omitempty" json:"name,omitempty"`
	Address           *string `xml:"Address,omitempty" json:"address,omitempty"`
	Address2          *string `xml:"Address_2,omitempty" json:"address_2,omitempty"`
	PostCode          *string `xml:"Post_Code,omitempty" json:"post_code,omitempty"`
	City              *string `xml:"City,omitempty" json:"city,omitempty"`
	County            *string `xml:"County,omitempty" json:"county,omitempty"`
	CountryRegionCode *string `xml:"Country_Region_Code,omitempty" json:"country_region_code,omitempty"`
	PhoneNo           *string `xml:"Phone_No,omitempty" json:"phone_no,omitempty"`
	// PrimaryContactNo      string `xml:"Primary_Contact_No,omitempty" json:"primary_contact_no,omitempty"`
	Contact               *string `xml:"Contact,omitempty" json:"contact,omitempty"`
	SearchName            *string `xml:"Search_Name,omitempty" json:"search_name,omitempty"`
	WeighbridgeSupplierID string  `xml:"Weighbridge_Supplier_ID,omitempty" json:"weighbridge_supplier_id,omitempty"`
	FaxNo                 *string `xml:"Fax_No,omitempty" json:"fax_no,omitempty"`
	Email                 *string `xml:"E_Mail,omitempty" json:"e_mail,omitempty"`
	// PayToVendorNo         string `xml:"Pay_to_Vendor_No,omitempty" json:"pay_to_vendor_no,omitempty"`
	GenBusPostingGroup string  `xml:"Gen_Bus_Posting_Group,omitempty" json:"gen_bus_posting_group,omitempty"`
	VATBusPostingGroup string  `xml:"VAT_Bus_Posting_Group,omitempty" json:"vat_bus_posting_group,omitempty"`
	VendorPostingGroup string  `xml:"Vendor_Posting_Group,omitempty" json:"vendor_posting_group,omitempty"`
	InvoiceDiscCode    *string `xml:"Invoice_Disc_Code,omitempty" json:"invoice_disc_code,omitempty"`
	ApplicationMethod  string  `xml:"Application_Method,omitempty" json:"application_method,omitempty"`
	ACN                string  `xml:"IRD_No,omitempty" json:"acn,omitempty"`
	ABN                *string `xml:"ABN,omitempty" json:"abn,omitempty"`
	ABNDivisionPartNo  string  `xml:"ABN_Division_Part_No,omitempty" json:"abn_division_part_no,omitempty"`
	Registered         bool    `xml:"Registered,omitempty" json:"registered,omitempty"`
}

type BackToCDSVendorResponse struct {
	VendorNo              string `json:"vendor_no,omitempty"`
	WeighbridgeSupplierID string `json:"weighbridge_supplier_id,omitempty"`
}

// For Reading Vendor key
type ReadResultVendor struct {
	XMLName xml.Name       `xml:"Envelope,omitempty" json:"envelope,omitempty"`
	Body    ReadBodyVendor `xml:"Body,omitempty" json:"body,omitempty"`
}

type ReadBodyVendor struct {
	ReadResult ReadResultDetailVendor `xml:"Read_Result,omitempty" json:"read_result,omitempty"`
}

type ReadResultDetailVendor struct {
	ReadVendor ReadWSVendor `xml:"WSVendor,omitempty" json:"ws_vendor,omitempty"`
}

type ReadWSVendor struct {
	Key string `xml:"Key,omitempty" json:"key,omitempty"`
	No  string `xml:"No,omitempty" json:"no,omitempty"`
}
