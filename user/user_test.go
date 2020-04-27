package user

import (
	"testing"

	"github.com/dgraph-io/badger"
	"github.com/zereraz/assistme/db"
	"github.com/zereraz/assistme/policy"
)

func InsertUserToDb(username string, t *testing.T) *User {
	user, err := NewUser("", username, 101, policy.DefaultPolicy)
	if err != nil {
		t.Error(err)
	}
	err = user.AddToDb()
	if err != nil {
		t.Errorf("Could not add user to db: %v", err)
	}
	return user

}

func TestUserInsert(t *testing.T) {
	username := "zereraz"
	user := InsertUserToDb(username, t)

	item, err := db.GetValueItem(user.GenerateKey())
	if err != nil {
		t.Error(err)
	}
	err = item.Value(func(val []byte) error {
		dbUser, err := ToUser(val)
		if err != nil {
			t.Error(err)
		}
		if !user.IsEqual(dbUser) {
			t.Error("Incorrect user unmarshaled")
		}
		return nil
	})
	if err != nil {
		t.Error(err)
	}
}

func TestDeleteUser(t *testing.T) {
	username := "zereraz"
	user := InsertUserToDb(username, t)

	userKey := user.GenerateKey()
	err := db.DeleteKey(userKey)
	if err != nil {
		t.Errorf("Could not delete user %v", err)
	}
	_, err = db.GetValueItem(userKey)
	if err != badger.ErrKeyNotFound {
		t.Error(err)
	}
}
