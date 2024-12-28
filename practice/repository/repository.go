package repository

import (
	"database/sql"
	"errors"
	"github.com/akmalfairuz/df-practice/practice/database"
)

var userRepo *User
var banRepo *Ban

func init() {
	userRepo = &User{db: database.Get()}
	banRepo = &Ban{db: database.Get()}
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
