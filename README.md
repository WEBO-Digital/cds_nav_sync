## Make .exe file
go build -o nav_sync_test.exe

## Run jobs by action
nav_sync_test.exe -action vendor_fetch 
nav_sync_test.exe -action vendor_sync
nav_sync_test.exe -action vendor_resync
nav_sync_test.exe -action invoice_fetch
nav_sync_test.exe -action invoice_sync
nav_sync_test.exe -action invoice_resync
nav_sync_test.exe -action ledger_entries_fetch
nav_sync_test.exe -action ledger_entries_sync
nav_sync_test.exe -action ledger_entries_resync

## Directory structure

data
 |_ vendor
   |_ pending
     |_ timestamp.json
   |_ done
     |_ timestamp.json



//How to sync data
pending
--- vendoers.json

done
--- vendoers.json


vendors
{
    "CDS-ID": {
        "hash"  : "dfKJdfjJdjf",
        "nav_id": "xdfjk"
    },
}

hash_records = {
    aaa: []
    bbb: []
    ccc: []
    ddd: 
}

loop: vendoers.json -> vender:
    cds_id = vender.cds_id
    hash  = md5(vender)
    hash_record = null

    if hash_records has cds_id
        hash_record = hash_records.cds_id

    if hash_record is null
        nav_id = nav.insertVender(vendor)
        
        if nav_id
            log success
            hash_records[cds_id] = {
                hash: hash
                nav_id: nav_id
            }
        else
            log failed

    if hash_record != null and hash_record.hash != hash
        /** @TODO: UPDATE THE VEDOR **/
    
    saveHashRecords(hash_records)