## Make .exe file
```
go build -o nav_sync_test.exe
```

## Run jobs by action
```
nav_sync_test.exe -action vendor_fetch 
nav_sync_test.exe -action vendor_sync
nav_sync_test.exe -action vendor_resync
nav_sync_test.exe -action invoice_fetch
nav_sync_test.exe -action invoice_sync
nav_sync_test.exe -action invoice_resync
nav_sync_test.exe -action ledger_entries_fetch
nav_sync_test.exe -action ledger_entries_sync
nav_sync_test.exe -action ledger_entries_resync
```


## Cron Job in NAV system
```
vendor_fetch
"C:\Users\rojan.shrestha\Desktop\nav syncing test\crons\vendor_fetch_cron.bat"
```
```
vendor_sync
"C:\Users\rojan.shrestha\Desktop\nav syncing test\crons\vendor_sync_cron.bat"
```
```
vendor_resync
"C:\Users\rojan.shrestha\Desktop\nav syncing test\crons\vendor_resync_cron.bat"
```
```
invoice_fetch
"C:\Users\rojan.shrestha\Desktop\nav syncing test\crons\invoice_fetch_cron.bat"
```
```
invoice_sync
"C:\Users\rojan.shrestha\Desktop\nav syncing test\crons\invoice_sync_cron.bat"
```
```
invoice_resync
"C:\Users\rojan.shrestha\Desktop\nav syncing test\crons\invoice_resync_cron.bat"
```
```
ledger_entries_fetch
"C:\Users\rojan.shrestha\Desktop\nav syncing test\crons\ledger_entries_fetch_cron.bat"
```
```
ledger_entries_sync
"C:\Users\rojan.shrestha\Desktop\nav syncing test\crons\ledger_entries_sync_cron.bat"
```
```
ledger_entries_resync
"C:\Users\rojan.shrestha\Desktop\nav syncing test\crons\ledger_entries_resync_cron.bat"
```


## Directory structure
```
data
 |_ vendor
   |_ pending
     |_ timestamp.json
   |_ done
     |_ timestamp.json
```

## How to sync data
```
pending
--- vendoers.json

done
--- vendoers.json
```

## Hash Function Example
```
vendors
{
    "CDS-ID": {
        "hash"  : "dfKJdfjJdjf",
        "nav_id": "xdfjk"
    },
}
```

## Algorithm to sync data
```
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

```

```
CGO_ENABLED=0 GOOS=linux go build -buildvcs=false
```

```sql
UPDATE refunds SET document_id=NULL, purchase_invoice_no=NULL
UPDATE customers SET nav_id=NULL
UPDATE payments SET nav_id=NULL
```
