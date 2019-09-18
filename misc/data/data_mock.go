// +build mock

package data

import (
	"fmt"
	"time"

	"eiko/misc/structures"
)

type Data struct {
	// Users users name inside the datastore
	Users string
}

// InitData return an initialised Data struct
func InitData(projID string) Data {
	var d Data
	d.Users = "Users"
	return d
}

// GetUser is used to find if a email is already used in the datastore
func (d Data) GetUser(UserMail string) (structures.User, error) {
	return structures.User{
		Email:     UserMail,
		Pass:      "$2a$10$OXJQl253CXXw.G/DADW3MO/bFhkuttZp5m4iwed83dDN4cZSD.hqe", // hashed password
		Created:   time.Now(),
		Validated: false,
	}, fmt.Errorf("test")
}

// StoreUser is used to store a user in the datastore
func (d Data) StoreUser(user structures.User) error {
	return nil
}