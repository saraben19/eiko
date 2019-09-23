// +build mock

package data

import (
	"fmt"
	"time"

	"eiko/misc/structures"
)

var (
	UserStored      bool
	Logged          bool
	GetUser         bool
	Inited          bool
	GetStore        bool
	StoreStore      bool
	StoreConsumable bool
	Error           error
	ErrTest         = fmt.Errorf("Test %s", "error")
	User            = structures.User{}
	UserTest        = structures.User{
		Email:     "test@test.ts",
		Pass:      "$2a$10$EVCZ/75E1TCgpOZFypJC4ejYDDTPk9lAGwLKGhp6jESMWfl/4Bl/e", // hashed password 'pass'
		Created:   time.Now(),
		Validated: false,
	}
	Store     = structures.Store{}
	StoreTest = structures.Store{
		Name:       "test store",
		Address:    "test store",
		Country:    "test store",
		Zip:        "test store",
		UserRating: 5,
	}
)

// Data container for all data relative variables
type Data struct {

	// User the user making the request. Got from the cookie in the header
	User structures.User
}

// InitData return an initialised Data struct
func InitData(projID string) Data {
	Inited = true
	var d Data
	return d
}

// GetUser is used to find if a email is already used in the datastore
func (d Data) GetUser(UserMail string) (structures.User, error) {
	GetUser = true
	return User, Error
}

// StoreUser is used to store a user in the datastore
func (d Data) StoreUser(user structures.User) error {
	UserStored = true
	return nil
}

// Log is used to store a log in the datastore
func (d Data) Log(user structures.Log) error {
	Logged = true
	return nil
}

// GetStore is used to find if a store is already in the datastore using
// it's name and location
func (d Data) GetStore(structures.Store) (structures.Store, error) {
	GetStore = true
	return Store, Error
}

// StoreStore is used to store a log in the datastore
func (d Data) StoreStore(store structures.Store) error {
	StoreStore = true
	return Error
}

// StoreConsumable is used to store a log in the datastore
func (d Data) StoreConsumable(consumable structures.Consumable) error {
	StoreConsumable = true
	return Error
}
