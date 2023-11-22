package invoice

import (
	"nav_sync/mods"
	filesystem "nav_sync/mods/afile_system"
	"nav_sync/mods/amanager"
)

func Fetch() {
	//Fetch vendor data
	response, err := amanager.Fetch(mods.INVOICE_FETCH_URL)
	if err != nil {
		mods.Console(err.Error())
	}
	mods.Console(response)

	//Save to pending file
	err = filesystem.Save(mods.INVOICE_PENDING_FILE_PATH, response)
	if err != nil {
		mods.Console(err.Error())
	}
	mods.Console("Successfully saved invoice to pending file")
}

func Sync() {
	//How to access pending directory
	//Case 1. Do this function run independently
	//Case 2. Do we run this function after fetch function is called

	//We go through Case 1.
	//Fetch all the pending files
	//Then sync one by one
	//After sync one file then move it to done folder.

	//Get All the vendor pending data
	fileNames, err := filesystem.GetAllFiles(mods.INVOICE_PENDING_FILE_PATH)
	if err != nil {
		mods.Console(err.Error())
	}

	//mods.Console(fileNames)

	if fileNames == nil || len(fileNames) < 1 {
		return
	}

	for i := 0; i < len(fileNames); i++ {
		//Sync vendor data to NAV
		//We assume here that Data are pushed to NAV

		//Move to done file
		err = filesystem.MoveFile(fileNames[i], mods.INVOICE_PENDING_FILE_PATH, mods.INVOICE_DONE_FILE_PATH)
		if err != nil {
			mods.Console(err.Error())
		}
	}
}
