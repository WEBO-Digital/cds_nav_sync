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
	RuncFunctionFromCommandArgument()
	//cmdrun.RuncFunctionFromCommandArgument()

	//Specify scheduler runner
	//cronjob.RunCron(1, 1, "vendor_fetch")

}

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
	case "payment_fetch":
		ledgerentries.Fetch()
	case "payment_sync":
		ledgerentries.Sync3()
	case "payment_resync":
		ledgerentries.ReSync()
	default:
		utils.Console("Invalid action!\n\nAvailable actions: vendor_fetch, vendor_sync, vendor_resync, invoice_fetch, invoice_sync, invoice_resync, payment_fetch, payment_sync, payment_resync")
		os.Exit(1)
	}
}

//To run
// nav_sync_test.exe -action vendor_fetch
