package main

import (
	"nav_sync/config"
	"nav_sync/mods/invoice"
	"nav_sync/mods/vendor"
)

func main() {
	//Load yaml config
	config.LoadYamlFile()

	//invoice
	// invoice.Fetch()
	// time.Sleep(10 * 1000)
	invoice.Sync()

	//vendor
	//vendor.Fetch()
	// time.Sleep(5 * 1000)
	vendor.Sync()

}
