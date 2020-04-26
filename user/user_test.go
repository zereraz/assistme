package user

import (
	"testing"

	"github.com/zereraz/assistme/db"
	"github.com/zereraz/assistme/policy"
)

func TestUserInsert(t *testing.T) {
	currDb, err := db.GetDb()
	username := "zereraz"
	user, err := NewUser("", username, 101, policy.DefaultPolicy)
	if err != nil {
		t.Error(err)
	}
	err = user.AddToDb(currDb)
	if err != nil {
		t.Errorf("Could not add user to db: %v", err)
	}

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
