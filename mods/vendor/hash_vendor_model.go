package vendor

//Here key string is CustomerNo -> Weighbridge_Supplier_ID
type HashVendorModel map[string]HashVendorEntry

type HashVendorEntry struct {
	Hash  string  `json:"hash"`
	NavID *string `json:"nav_id"`
}
