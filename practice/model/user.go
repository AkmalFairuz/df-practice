package model

type User struct {
	ID          int    `db:"id"`
	DisplayName string `db:"displayName"`
	Name        string `db:"name"`
	XUID        string `db:"xuid"`

	RankName     string `db:"rankName"`
	RankExpireAt int64  `db:"rankExpireAt"`

	Exp int64 `db:"exp"`

	RegisteredAt int64 `db:"registeredAt"`
	LastSeenAt   int64 `db:"lastSeenAt"`
}

type CreateUser struct {
	DisplayName string `db:"displayName"`
	XUID        string `db:"xuid"`
}
