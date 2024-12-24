package ffa

import (
	"github.com/akmalfairuz/df-practice/practice/user"
	"sync"
	"sync/atomic"
	"time"
)

type Participant struct {
	u *user.User

	lastAttackedMu     sync.Mutex
	lastAttackedByXUID string
	lastAttackedAt     time.Time

	combatTimer atomic.Int32

	kills      atomic.Int32
	killStreak atomic.Int32
	deaths     atomic.Int32

	lastSpawn time.Time
}

func (par *Participant) StoreLastAttackedBy(xuid string) {
	par.lastAttackedMu.Lock()
	defer par.lastAttackedMu.Unlock()
	par.lastAttackedByXUID = xuid
	par.lastAttackedAt = time.Now()
}

func (par *Participant) LastAttackedByWithMaxDuration(maxDuration time.Duration) string {
	par.lastAttackedMu.Lock()
	defer par.lastAttackedMu.Unlock()
	if time.Since(par.lastAttackedAt) > maxDuration {
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
