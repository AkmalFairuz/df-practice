package repository

import (
	"github.com/akmalfairuz/df-practice/practice/model"
	"github.com/jmoiron/sqlx"
)

type Ban struct {
	db *sqlx.DB
}

func NewBan(db *sqlx.DB) *Ban {
	return &Ban{db: db}
}

func (r *Ban) FindByPlayerID(playerId int) (model.Ban, error) {
	var b model.Ban
	err := r.db.Get(&b, "SELECT * FROM bans WHERE playerId = ?", playerId)
	return b, err
}

func (r *Ban) Create(ban model.Ban) (int64, error) {
	ret, err := r.db.Exec("INSERT INTO bans (playerId, reason, expireAt, createdAt) VALUES (?, ?, ?, ?)", ban.PlayerID, ban.Reason, ban.ExpireAt, ban.CreatedAt)
	if err != nil {
		return 0, err
	}
	return ret.LastInsertId()
}

func (r *Ban) DeleteByPlayerID(playerId int) (int64, error) {
	res, err := r.db.Exec("DELETE FROM bans WHERE playerId = ?", playerId)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}
