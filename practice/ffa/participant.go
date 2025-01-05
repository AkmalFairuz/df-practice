package ffa

import (
	"github.com/akmalfairuz/df-practice/practice/user"
	"github.com/df-mc/atomic"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"sync"
	"time"
)

type Participant struct {
	u *user.User

	lastAttackedMu     sync.Mutex
	lastAttackedByXUID string
	lastAttackedAt     atomic.Value[time.Time]

	combatTimer atomic.Int32

	kills      atomic.Int32
	killStreak atomic.Int32
	deaths     atomic.Int32

	lastPearlThrow atomic.Value[time.Time]

	lastSpawn atomic.Value[time.Time]
}

func (par *Participant) StoreLastAttackedBy(xuid string) {
	par.lastAttackedMu.Lock()
	defer par.lastAttackedMu.Unlock()
	par.lastAttackedByXUID = xuid
	par.lastAttackedAt.Store(time.Now())
}

func (par *Participant) LastAttackedByWithMaxDuration(maxDuration time.Duration) string {
	par.lastAttackedMu.Lock()
	defer par.lastAttackedMu.Unlock()
	if time.Since(par.lastAttackedAt.Load()) > maxDuration {
		return ""
	}
	return par.lastAttackedByXUID
}

func (par *Participant) LastAttackedBy() string {
	return par.LastAttackedByWithMaxDuration(8 * time.Second)
}

func (par *Participant) InCombat() bool {
	return par.combatTimer.Load() > 0
}

func (par *Participant) Combat() int {
	return int(par.combatTimer.Load())
}

func (par *Participant) Player(tx *world.Tx) (*player.Player, bool) {
	e, ok := par.u.EntityHandle().Entity(tx)
	if !ok {
		return nil, false
	}
	return e.(*player.Player), true
}

func (par *Participant) OnKill() {
	par.kills.Inc()
	par.killStreak.Inc()
}
