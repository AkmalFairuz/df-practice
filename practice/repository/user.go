package repository

import (
	"github.com/akmalfairuz/df-practice/practice/model"
	"github.com/jmoiron/sqlx"
	"strings"
	"time"
)

type User struct {
	db *sqlx.DB
}

func NewUser(db *sqlx.DB) *User {
	return &User{db: db}
}

func (r *User) FindByID(id int) (model.User, error) {
	var u model.User
	err := r.db.Get(&u, "SELECT * FROM users WHERE id = ?", id)
	return u, err
}

func (r *User) FindByXUID(xuid string) (model.User, error) {
	var u model.User
	err := r.db.Get(&u, "SELECT * FROM users WHERE xuid = ?", xuid)
	return u, err
}

func (r *User) FindByName(name string) (model.User, error) {
	var u model.User
	err := r.db.Get(&u, "SELECT * FROM users WHERE name = ?", strings.ToLower(name))
	return u, err
}

func (r *User) Create(create model.CreateUser) (int64, error) {
	ret, err := r.db.Exec("INSERT INTO users (name, displayName, xuid, lastSeenAt, registeredAt) VALUES (?, ?, ?, ?, ?)", strings.ToLower(create.DisplayName), create.DisplayName, create.XUID, time.Now().Unix(), time.Now().Unix())
	if err != nil {
		return 0, err
	}
	return ret.LastInsertId()
}

func (r *User) SetDisplayName(id int, displayName string) error {
	_, err := r.db.Exec("UPDATE users SET displayName = ?, name = ? WHERE id = ?", displayName, strings.ToLower(displayName), id)
	return err
}

func (r *User) SynchronizeLastSeen(id int) error {
	_, err := r.db.Exec("UPDATE users SET lastSeenAt = ? WHERE id = ?", time.Now().Unix(), id)
	return err
}
