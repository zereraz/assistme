package user

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/dgraph-io/badger"
	"github.com/teris-io/shortid"
	"github.com/zereraz/assistme/config"
	"github.com/zereraz/assistme/policy"
	"github.com/zereraz/assistme/statistics"
	"github.com/zereraz/assistme/utils"
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
	return utils.GenerateHash(u.Username)
}

func (u *User) GenerateKey() []byte {
	return []byte(fmt.Sprintf("user%s%d", config.KeyDelim, u.generateHash()))
}

func (u *User) AddToDb(db *badger.DB) error {
	userJson, err := json.Marshal(u)
	if err != nil {
		return err
	}
	err = db.Update(func(txn *badger.Txn) error {
		err := txn.Set(u.GenerateKey(), userJson)
		return err
	})
	return err
}

func (u *User) IsEqual(newUser *User) bool {
	return u.Name == newUser.Name &&
		u.Username == newUser.Username &&
		u.Id == newUser.Id &&
		u.ChannelId == newUser.ChannelId &&
		*u.Policy == *newUser.Policy &&
		*u.Statistics == *newUser.Statistics
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

func ToUser(marshaledUser []byte) (*User, error) {
	user := &User{}
	err := json.Unmarshal(marshaledUser, user)
	if err != nil {
		return nil, err
	}
	return user, nil
}
