package repository

import (
	"database/sql"
	"errors"
	"github.com/akmalfairuz/df-practice/practice/database"
)

var userRepo *User
var banRepo *Ban

func init() {
	userRepo = NewUser(database.Get())
	banRepo = NewBan(database.Get())
}

func UserRepo() *User {
	return userRepo
}

func BanRepo() *Ban {
	return banRepo
}

func IsNotExists(err error) bool {
	return errors.Is(err, sql.ErrNoRows)
}
