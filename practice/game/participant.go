package game

import "github.com/akmalfairuz/df-practice/practice/user"

type Participant struct {
	xuid  string
	state ParticipantState
}

func (p *Participant) XUID() string {
	return p.xuid
}

func (p *Participant) User() *user.User {
	return user.GetByXUID(p.xuid)
}

func (p *Participant) IsPlaying() bool {
	return p.state == ParticipantStatePlaying
}

func (p *Participant) IsSpectating() bool {
	return p.state == ParticipantStateSpectating
}
