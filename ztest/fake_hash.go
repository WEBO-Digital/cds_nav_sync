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
		Hash: "dmkdfkdf",
	})

	fmt.Println(recs.GetHash("nn1"))

	recs.Save()
}
