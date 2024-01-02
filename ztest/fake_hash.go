package main

import (
	"fmt"
	"nav_sync/mods/hashrecs"
)

func main() {
	recs := hashrecs.HashRecs{
		Name: "tee",
	}

	recs.Load()

	data := `
	{
		"abn": null,
		"acn": "474396551",
		"address": "QLD",
		"address_2": "Glass House",
		"application_method": "Manual",
		"city": "Cityville",
		"contact": "John Doe",
		"country_region_code": "AU",
		"county": "Countyshire",
		"email": "test@me.com",
		"fax_no": "0491570156",
		"gen_bus_posting_group": "SPL",
		"name": "N2f-test",
		"pay_to_vendor_no": "TEST 5",
		"phone_no": "0296956888",
		"post_code": null,
		"primary_contact_no": "0491570156",
		"registered": false,
		"vat_bus_posting_group": "LOCAL",
		"vendor_posting_group": "DOMESTIC",
		"weighbridge_supplier_id": "C14"
	}
	`

	recs.Set("nn1", hashrecs.HashRec{
		Hash: hashrecs.Hash(data), //"dmkdfkdf",
	})

	fmt.Println("hash------------>", recs.GetHash("nn1"))

	recs.Save()
}
