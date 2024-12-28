package game

import (
	"github.com/akmalfairuz/df-practice/practice/user"
	"github.com/df-mc/atomic"
	"time"
)

type Participant struct {
	xuid  string
	state ParticipantState

	pearlCooldown atomic.Value[time.Time]
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

func (p *Participant) PearlCooldown() time.Time {
	return p.pearlCooldown.Load()
}

func (p *Participant) SetPearlCooldown(d time.Time) {
	p.pearlCooldown.Store(d)
}
