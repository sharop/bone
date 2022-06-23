package log

import store2 "github.com/sharop/nopaldb/store"

type OTx struct {
	store    *store2.Database
	entries  []*Entry
	metadata TxMetadata // Meta for transaction
	closed   bool
}

type Entry struct {
	Key      []byte
	Metadata FMetadata //Meta for entry. This meta should be difined in a code table. Now Just a String
	Value    []byte
}
