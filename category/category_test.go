package category

import (
	"testing"

	"github.com/zereraz/assistme/db"
	"github.com/zereraz/assistme/policy"
	"github.com/zereraz/assistme/user"
)

func TestCategoryInsert(t *testing.T) {
	currDb, err := db.GetDb()
	username := "zereraz"
	user, err := user.NewUser("", username, 101, policy.DefaultPolicy)
	if err != nil {
		t.Error(err)
	}

	err = user.AddToDb(currDb)
	if err != nil {
		t.Errorf("Could not add user to db: %v", err)
	}
	userKey := user.GenerateKey()
	categoryName := "books"
	category, err := NewCategory(categoryName, "", username)
	if err != nil {
		t.Error(err)
	}

	err = category.AddToDb(currDb, userKey)
	if err != nil {
		t.Error(err)
	}

	item, err := db.GetValueItem(category.GenerateKey(userKey))
	if err != nil {
		t.Error(err)
	}
	err = item.Value(func(val []byte) error {
		dbCategory, err := ToCategory(val)
		if err != nil {
			t.Error(err)
		}
		if !category.IsEqual(dbCategory) {
			t.Error("Incorrect category unmarshaled")
		}
		return nil
	})
	if err != nil {
		t.Error(err)
	}
}
