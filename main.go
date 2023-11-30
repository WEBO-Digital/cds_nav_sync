package main

import (
	"nav_sync/config"
	"nav_sync/mods/invoice"
	ledgerentries "nav_sync/mods/ledger_entries"
	"nav_sync/mods/vendor"
)

func main() {
	//Load yaml config
	config.LoadYamlFile()

	//invoice
	//invoice.Fetch()
	invoice.Sync()

	//vendor
	//vendor.Fetch()
	vendor.Sync()

	//ledger_entries
	//ledgerentries.Fetch()
	ledgerentries.Sync()
}
