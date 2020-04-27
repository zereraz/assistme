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

// get item given value
// will return ErrKeyNotFound
// if key does not exist
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

// get value given key
func GetValue(key []byte) ([]byte, error) {
	item, err := GetValueItem(key)
	if err != nil {
		return nil, err
	}
	val, err := GetValueFromItem(item)
	if err != nil {
		return nil, err
	}
	return val, nil
}

// given an item get value from it
func GetValueFromItem(item *badger.Item) ([]byte, error) {
	var valCopy []byte
	err := item.Value(func(val []byte) error {
		valCopy = append([]byte{}, val...)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return valCopy, nil
}

// Delete value by key
func DeleteKey(key []byte) error {
	currDb, err := GetDb()
	if err != nil {
		return err
	}
	return currDb.Update(func(txn *badger.Txn) error {
		return txn.Delete(key)
	})
}
