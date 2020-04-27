package category

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/dgraph-io/badger"
	"github.com/teris-io/shortid"
	"github.com/zereraz/assistme/config"
	"github.com/zereraz/assistme/utils"
)

const (
	Namespace = "category"

	// categories that get made per user
)

var (
	DefaultCategoryNames = []string{"Files", "Tasks", "Ideas", "Thoughts", "Notes"}
)

type Category struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Username    string `json:"username"`
}

// create new category
func NewCategory(name, description, username string) (*Category, error) {
	if name == "" {
		return nil, errors.New("Cannot create category with empty name")
	}
	categoryId, err := shortid.Generate()
	if err != nil {
		return nil, err
	}
	return &Category{categoryId, name, description, username}, nil
}

func (c *Category) generateHash() uint64 {
	return utils.GenerateHash(c.Name)
}

// generate category key
func (c *Category) GenerateKey(userKey []byte) []byte {
	return append(userKey, []byte(fmt.Sprintf("%s%s%s%d", config.KeyDelim, Namespace, config.KeyDelim, c.generateHash()))...)
}

// Add category to db
func (c *Category) AddToDb(db *badger.DB, userKey []byte) error {
	categoryJson, err := json.Marshal(c)
	if err != nil {
		return err
	}
	err = db.Update(func(txn *badger.Txn) error {
		err := txn.Set(c.GenerateKey(userKey), categoryJson)
		return err
	})
	return err
}

// compare if categories equal
func (c *Category) IsEqual(newCategory *Category) bool {
	return c.Name == newCategory.Name &&
		c.Id == newCategory.Id &&
		c.Description == newCategory.Description &&
		c.Username == newCategory.Username
}

// stored marshaled value to Category
func ToCategory(marshaledCategory []byte) (*Category, error) {
	category := &Category{}
	err := json.Unmarshal(marshaledCategory, category)
	if err != nil {
		return nil, err
	}
	return category, nil
}
