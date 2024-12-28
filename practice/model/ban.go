package model

import "time"

type Ban struct {
	ID        int    `db:"id"`
	PlayerID  int    `db:"playerId"`
	Reason    string `db:"reason"`
	ExpireAt  int64  `db:"expireAt"`
	CreatedAt int64  `db:"createdAt"`
}

func (b Ban) Remaining() (days int, hours int, minutes int) {
	now := time.Now().Unix()
	days = int((b.ExpireAt - now) / 86400)
	hours = int((b.ExpireAt - now) % 86400 / 3600)
	minutes = int((b.ExpireAt - now) % 3600 / 60)
	return
}
