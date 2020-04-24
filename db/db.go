package db

import (
	"github.com/dgraph-io/badger"
	"github.com/raunaqrox/assistme/config"
	"github.com/raunaqrox/assistme/log"
)

var Db *badger.DB

func SetupDb() error {
	db, err := badger.Open(badger.DefaultOptions(config.DbPath))
	if err != nil {
		return err
	}
	Db = db
	log.Log.Println("Db is setup")
	return nil
}

func Cleanup() {
	if Db != nil {
		Db.Close()
	}
}
