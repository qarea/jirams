// Package main provides main service.
package main

import (
	"github.com/boltdb/bolt"
	"github.com/qarea/jirams/cfg"

	"github.com/powerman/narada-go/narada"
)

var l = narada.NewLog("")

func main() {
	db, err := bolt.Open("var/bolt/store.db", 0666, nil)
	if err != nil {
		l.Fatal(err)
	}
	defer func() { _ = db.Close() }()

	start(appParams{
		BoltDB:    db,
		PublicKey: cfg.RSAPublicKey,
	})
}
