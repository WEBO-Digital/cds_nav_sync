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