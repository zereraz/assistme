package db

import (
	"github.com/dgraph-io/badger"
	"github.com/zereraz/assistme/config"
	"github.com/zereraz/assistme/log"
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

func GetValueItem(key []byte) (*badger.Item, error) {
	var item *badger.Item
	db, err := GetDb()
	if err != nil {
		return nil, err
	}
	err = db.View(func(txn *badger.Txn) error {
		item, err = txn.Get(key)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return item, nil
}

func Cleanup() {
	if Db != nil {
		Db.Close()
	}
}
