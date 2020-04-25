package user

import (
	"testing"

	"github.com/raunaqrox/assistme/db"
	"github.com/raunaqrox/assistme/policy"
)

func TestUserInsert(t *testing.T) {
	db, err := db.GetDb()
	username := "zereraz"
	user, err := NewUser("", username, 101, policy.DefaultPolicy)
	if err != nil {
		t.Errorf("%v", err)
	}
	err = user.AddToDb(db)
	if err != nil {
		t.Errorf("Could not add user to db: %v", err)
	}
}
