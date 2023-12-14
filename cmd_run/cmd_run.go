package cmdrun

import (
	"flag"
	"nav_sync/mods/invoice"
	ledgerentries "nav_sync/mods/ledger_entries"
	"nav_sync/mods/vendor"
	"nav_sync/utils"
	"os"
)

func RuncFunctionFromCommandArgument() {
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
	case "vendor_resync":
		vendor.ReSync()
	case "invoice_fetch":
		invoice.Fetch()
	case "invoice_sync":
		invoice.Sync3()
	case "invoice_resync":
		invoice.ReSync()
	case "ledgerentries_fetch":
		ledgerentries.Fetch()
	case "ledgerentries_sync":
		ledgerentries.Sync3()
	case "ledgerentries_resync":
		ledgerentries.ReSync()
	default:
		utils.Console("Invalid action. Available actions: vendor_fetch, vendor_sync, vendor_resync, invoice_fetch, invoice_sync, invoice_resync, ledgerentries_fetch, ledgerentries_sync, ledgerentries_resync")
		os.Exit(1)
	}
}
