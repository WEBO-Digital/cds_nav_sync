package main

import (
	"nav_sync/mods/vendor"
)

func main() {
	//vendor
	// vendor.Fetch()
	// time.Sleep(5 * 1000)
	vendor.Sync()

	//invoice
	// invoice.Fetch()
	// time.Sleep(10 * 1000)
	// invoice.Sync()
}
