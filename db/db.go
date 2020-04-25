package db

import (
	"github.com/dgraph-io/badger"
	"github.com/raunaqrox/assistme/config"
	"github.com/raunaqrox/assistme/log"
)

var Db *badger.DB

func SetupDb() (*badger.DB, error) {
	db, err := badger.Open(badger.DefaultOptions(config.DbPath))
	if err != nil {
		return nil, err
	}
	Db = db
	log.Log.Println("Db is setup")
	return db, nil
}

func GetDb() (*badger.DB, error) {
	if Db != nil {
		return Db, nil
	}
	return SetupDb()
}

func Cleanup() {
	if Db != nil {
		Db.Close()
	}
}
