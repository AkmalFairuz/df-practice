package igame

import "github.com/akmalfairuz/df-practice/practice/user"

type IParticipant interface {
	XUID() string
	User() *user.User
	IsSpectating() bool
	IsPlaying() bool
}
