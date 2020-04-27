package user

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/dgraph-io/badger"
	"github.com/teris-io/shortid"
	"github.com/zereraz/assistme/category"
	"github.com/zereraz/assistme/config"
	"github.com/zereraz/assistme/db"
	"github.com/zereraz/assistme/policy"
	"github.com/zereraz/assistme/statistics"
	"github.com/zereraz/assistme/utils"
)

const (
	Namespace = "user"
)

// Errors
var (
	ErrUserNotFound = errors.New("User does not exist")
)

type User struct {
	Name       string                 `json:"name"`
	Username   string                 `json:"username"`
	Id         string                 `json:"id"`
	ChannelId  int64                  `json:"channelId"`
	Policy     *policy.Policy         `json:"policy"`
	Statistics *statistics.Statistics `json:"statistics"`
	Categories []*category.Category   `json:"categories"`
}

// User hash
func (u *User) generateHash() uint64 {
	return utils.GenerateHash(u.Username)
}

// Generate key for user
func (u *User) GenerateKey() []byte {
	return []byte(fmt.Sprintf("%s%s%d", Namespace, config.KeyDelim, u.generateHash()))
}

// Add user to Db
func (u *User) AddToDb() error {
	currDb, err := db.GetDb()
	if err != nil {
		return err
	}
	userJson, err := json.Marshal(u)
	if err != nil {
		return err
	}
	err = currDb.Update(func(txn *badger.Txn) error {
		err := txn.Set(u.GenerateKey(), userJson)
		return err
	})
	return err
}

// Compare if users equal
func (u *User) IsEqual(newUser *User) bool {
	return u.Name == newUser.Name &&
		u.Username == newUser.Username &&
		u.Id == newUser.Id &&
		u.ChannelId == newUser.ChannelId &&
		*u.Policy == *newUser.Policy &&
		*u.Statistics == *newUser.Statistics
}

// create new user with empty statistics and category
// username mandatory
func NewUser(name, username string, channelId int64, userPolicy *policy.Policy) (*User, error) {
	userId, err := shortid.Generate()
	if err != nil {
		return nil, err
	}
	if userPolicy == nil {
		userPolicy = policy.DefaultPolicy
	}
	if username == "" {
		return nil, errors.New("username cannot be empty")
	}
	var categories []*category.Category
	for _, name := range category.DefaultCategoryNames {
		category, err := category.NewCategory(name, "", username)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	return &User{name, username, userId, channelId, userPolicy, statistics.NewStatistics(), categories}, nil
}

// stored marshaled value to User
func ToUser(marshaledUser []byte) (*User, error) {
	user := &User{}
	err := json.Unmarshal(marshaledUser, user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// Given username fetch and get User object
func FetchUser(username string) (*User, error) {
	userKey := []byte(fmt.Sprintf("%s%s%d", Namespace, config.KeyDelim, utils.GenerateHash(username)))
	marshaledUser, err := db.GetValue(userKey)

	if err == badger.ErrKeyNotFound {
		return nil, ErrUserNotFound
	}
	// any other error
	if err != nil {
		return nil, err
	}

	user, err := ToUser(marshaledUser)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// Add category to user given username
func UpdateCategoryToUser(username string, c *category.Category) error {
	u, err := FetchUser(username)
	if err != nil {
		return err
	}
	u.Categories = append(u.Categories, c)
	return u.AddToDb()
}

func (u *User) DeleteUser() error {
	return db.DeleteKey(u.GenerateKey())
}

// TODO:
// Add method to delete by just username
// Add method to generate key just by username
