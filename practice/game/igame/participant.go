package igame

import (
	"github.com/akmalfairuz/df-practice/practice/user"
	"time"
)

type IParticipant interface {
	XUID() string
	User() *user.User
	IsSpectating() bool
	IsPlaying() bool
	PearlCooldown() time.Time
	SetPearlCooldown(time.Time)
}
