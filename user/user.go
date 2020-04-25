package user

import (
	"encoding/json"
	"errors"
	"fmt"
	"hash/fnv"

	"github.com/dgraph-io/badger"
	"github.com/raunaqrox/assistme/policy"
	"github.com/raunaqrox/assistme/statistics"
	"github.com/teris-io/shortid"
)

type User struct {
	Name       string                 `json:"name"`
	Username   string                 `json:"username"`
	Id         string                 `json:"id"`
	ChannelId  int64                  `json:"channelId"`
	Policy     *policy.Policy         `json:"policy"`
	Statistics *statistics.Statistics `json:"statistics"`
}

func (u *User) generateHash() uint64 {
	hash := fnv.New64a()
	hash.Write([]byte(u.Username))
	return hash.Sum64()
}

func (u *User) generateKey() []byte {
	return []byte(fmt.Sprintf("user:%d", u.generateHash()))
}

func (u *User) AddToDb(db *badger.DB) error {
	userJson, err := json.Marshal(u)
	if err != nil {
		return err
	}
	err = db.Update(func(txn *badger.Txn) error {
		err := txn.Set(u.generateKey(), userJson)
		return err
	})
	return err
}

func NewUser(name, username string, channelId int64, policy *policy.Policy) (*User, error) {
	userId, err := shortid.Generate()
	if err != nil {
		return nil, err
	}
	if username == "" {
		return nil, errors.New("username cannot be empty")
	}
	return &User{name, username, userId, channelId, policy, statistics.NewStatistics()}, nil
}
