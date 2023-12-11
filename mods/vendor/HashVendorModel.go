package vendor

type HashVendorModel map[string]HashVendorEntry

type HashVendorEntry struct {
	Hash  string  `json:"hash"`
	NavID *string `json:"nav_id"`
}
