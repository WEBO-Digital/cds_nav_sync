package utils

//File Paths

// 1. Vendor
const VENDOR_PENDING_FILE_PATH = "/data/vendor/pending/"
const VENDOR_DONE_FILE_PATH = "/data/vendor/processed/"
const VENDOR_LOG_PATH = "/data/vendor/logs/"

// c. Hash
const VENDOR_HASH_FILE_PATH = "/data/vendor/hashrecs/"
const VENDOR_HASH_DB = "vendor.hash"

// 2. Invoice
const INVOICE_PENDING_FILE_PATH = "/data/invoice/pending/"
const INVOICE_DONE_FILE_PATH = "/data/invoice/processed/"
const INVOICE_LOG_PATH = "/data/invoice/logs/"

// c. Hash
const INVOICE_HASH_FILE_PATH = "/data/invoice/hashrecs/"
const INVOICE_HASH_DB = "invoice.hash"

// 3. Ledger Entries
// a. Pending
const LEDGER_ENTRIES_PENDING_FILE_PATH = "/data/ledger_entries/pending/"
const LEDGER_ENTRIES_PENDING_LOG_FILE_PATH = "/data/alogs/ledger_entries/pending/"
const LEDGER_ENTRIES_PENDING_FAILURE = "done.failure.log"
const LEDGER_ENTRIES_PENDING_SUCCESS = "done.success.log"

// b. Done
const LEDGER_ENTRIES_DONE_FILE_PATH = "/data/ledger_entries/done/"
const LEDGER_ENTRIES_DONE_LOG_FILE_PATH = "/data/alogs/ledger_entries/done/"
const LEDGER_ENTRIES_DONE_FAILURE = "done.failure.log"
const LEDGER_ENTRIES_DONE_SUCCESS = "done.success.log"

// c. Hash
const LEDGER_ENTRIES_HASH_FILE_PATH = "/data/ledger_entries/hashrecs/"
const LEDGER_ENTRIES_HASH_DB = "ledger_entries.hash"
