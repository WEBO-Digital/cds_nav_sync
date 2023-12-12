package main

import (
	"fmt"
	"nav_sync/mods/hashrecs"
)

func main() {
	recs := hashrecs.HashRecs{
		Name: "tee",
	}

	recs.Load()

	recs.Set("nn1", hashrecs.HashRec{
		Hash: hashrecs.Hash("dmkdfkdf"), //"dmkdfkdf",
	})

	fmt.Println("hash------------>", recs.GetHash("nn1"))

	recs.Save()
}
