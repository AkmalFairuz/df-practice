package igame

import (
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/session"
	"github.com/df-mc/dragonfly/server/world"
	"time"
)

type User interface {
	XUID() string
	Ping() int
	EntityHandle() *world.EntityHandle
	SendScoreboard(lines []string)
	Player(tx *world.Tx) (*player.Player, bool)
	Conn() session.Conn
	Translatef(format string, args ...any) string
	Messaget(translationName string, args ...any)
}

type IParticipant interface {
	XUID() string
	User() User
	IsSpectating() bool
	IsPlaying() bool
	PearlCooldown() time.Time
	SetPearlCooldown(time.Time)
}
