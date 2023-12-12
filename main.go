package main

import (
	"flag"
	"nav_sync/config"
	"nav_sync/mods/invoice"
	ledgerentries "nav_sync/mods/ledger_entries"
	"nav_sync/mods/vendor"
	"nav_sync/utils"
	"os"
)

func main() {
	//Load yaml config
	config.LoadYamlFile()
	runcFunctionFromCommandArgument()
}

func runcFunctionFromCommandArgument() {
	//Define a command-line flag named "action" with a default value of "defaultAction"
	action := flag.String("action", "defaultAction", "Specify the action to perform")

	//Parse the command-line arguments
	flag.Parse()

	//Call the appropratiate function based on the provided action
	switch *action {
	case "vendor_fetch":
		vendor.Fetch()
	case "vendor_sync":
		vendor.Sync3()
	case "invoice_fetch":
		invoice.Fetch()
	case "invoice_sync":
		invoice.Sync3()
	case "ledgerentries_fetch":
		ledgerentries.Fetch()
	case "ledgerentries_sync":
		ledgerentries.Sync3()
	default:
		utils.Console("Invalid action. Available actions: invoice_fetch, invoice_sync, vendor_fetch, vendor_sync, ledgerentries_fetch, ledgerentries_sync")
		os.Exit(1)
	}
}

//To run
// nav_sync_test.exe -action vendor_fetch
