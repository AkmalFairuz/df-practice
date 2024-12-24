package ffa

import (
	"errors"
	"github.com/akmalfairuz/df-practice/practice/helper"
	"github.com/akmalfairuz/df-practice/practice/user"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/entity"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"math/rand"
	"sync"
	"time"
)

type Arena struct {
	w *world.World

	u  map[string]*Participant
	mu sync.RWMutex

	spawns []helper.Location

	dropAllowed bool
	onSendKit   func(*player.Player) error

	icon string

	allowBuild bool

	placedBlocksMu sync.RWMutex
	placedBlocks   map[cube.Pos]world.Block
}

func New(w *world.World) *Arena {
	return &Arena{
		w: w,
		u: map[string]*Participant{},
	}
}

func (a *Arena) Participants() map[string]*Participant {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.u
}

func (a *Arena) ParticipantByXUID(xuid string) (*Participant, bool) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	p, ok := a.u[xuid]
	return p, ok
}

func (a *Arena) IsInArena(u *user.User) bool {
	a.mu.RLock()
	defer a.mu.RUnlock()
	_, ok := a.u[u.XUID()]
	return ok
}

func (a *Arena) Init() error {
	a.w.SetTime(3000)
	a.w.StopTime()
	a.w.StopRaining()
	a.w.StopThundering()
	a.w.StopWeatherCycle()

	go a.startTicking()
	return nil
}

func (a *Arena) startTicking() {
	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()

	currentTick := int64(0)

	for {
		select {
		case <-ticker.C:
			currentTick++
			a.handleTick(currentTick)
		}
	}
}

func (a *Arena) handleTick(currentTick int64) {
	<-a.w.Exec(func(tx *world.Tx) {
		for _, par := range a.Participants() {
			if currentTick%20 == 0 {
				par.combatTimer.Store(max(0, par.combatTimer.Load()-1))
				a.sendUserScoreboard(par, tx)
			}
		}
	})
}

func (a *Arena) sendUserScoreboard(p *Participant, tx *world.Tx) {
	combatInfo := "<empty>"
	if p.combatTimer.Load() > 0 {
		combatInfo = p.u.Translatef("ffa.scoreboard.combat.timer", p.combatTimer.Load())
	}

	p.u.SendScoreboard(p.u.Translatef("scoreboard.title"), []string{
		"",
		p.u.Translatef("ffa.scoreboard.your.kills", p.kills.Load()),
		p.u.Translatef("ffa.scoreboard.your.deaths", p.deaths.Load()),
		p.u.Translatef("ffa.scoreboard.your.streak", p.killStreak.Load()),
		"",
		p.u.Translatef("ffa.scoreboard.players", len(a.u)),
		p.u.Translatef("scoreboard.your.ping", p.u.Session().Latency().Milliseconds()),
		combatInfo,
		"",
		p.u.Translatef("scoreboard.footer"),
	})
}

func (a *Arena) Join(p *player.Player, tx *world.Tx) error {
	u := user.Get(p)

	if u.CurrentFFAArena() != nil {
		return errors.New("user already in arena")
	}

	a.mu.RLock()
	if _, ok := a.u[u.XUID()]; ok {
		a.mu.RUnlock()
		return errors.New("user already in this arena, this should not happen after u.CurrentFFAArena() check before")
	}
	par := &Participant{u: u}
	a.u[u.XUID()] = par
	a.mu.RUnlock()

	u.SetCurrentFFAArena(a)

	selectedSpawn := a.spawns[rand.Intn(len(a.spawns))]

	tx.RemoveEntity(p)

	<-a.w.Exec(func(tx2 *world.Tx) {
		newP := tx2.AddEntity(p.H()).(*player.Player)
		selectedSpawn.TeleportPlayer(newP)

		helper.ResetPlayer(newP)
		_ = a.sendKit(newP)
		newP.SetGameMode(world.GameModeSurvival)
		helper.UpdatePlayerNameTagWithHealth(newP, 0)
	})

	a.sendUserScoreboard(par, tx)

	return nil
}

func (a *Arena) Respawn(p *player.Player, tx *world.Tx) error {
	a.mu.Lock()
	par, ok := a.u[p.XUID()]
	a.mu.Unlock()

	if !ok {
		return errors.New("user is not in this arena")
	}

	helper.ResetPlayer(p)
	selectedSpawn := a.spawns[rand.Intn(len(a.spawns))]
	selectedSpawn.TeleportPlayer(p)
	par.combatTimer.Store(0)
	par.lastAttackedAt = time.Unix(0, 0)
	helper.UpdatePlayerNameTagWithHealth(p, 0)
	_ = a.sendKit(p)
	return nil
}

func (a *Arena) sendKit(p *player.Player) error {
	if a.onSendKit == nil {
		return nil
	}
	return (a.onSendKit)(p)
}

// Quit removes a player from the arena. Caller should teleport the player to lobby.
func (a *Arena) Quit(p *player.Player) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	u := user.Get(p)

	if _, ok := a.u[u.XUID()]; ok {
		return errors.New("user is not in this arena")
	}

	delete(a.u, u.XUID())
	u.SetCurrentFFAArena(nil)
	return nil
}

func (a *Arena) DropAllowed() bool {
	return a.dropAllowed
}

func (a *Arena) BroadcastMessaget(translationName string, args ...any) {
	for _, p := range a.u {
		p.u.Messaget(translationName, args...)
	}
}

func (a *Arena) HandleHurt(ctx *player.Context, damage *float64, immune bool, immunity *time.Duration, src world.DamageSource) {
	par, ok := a.u[ctx.Val().XUID()]
	if !ok {
		return
	}

	death := *damage >= ctx.Val().Health()

	if death {
		ctx.Cancel()
	}

	handledDeathMessage := false

	switch cause := src.(type) {
	case *entity.AttackDamageSource:
		attacker, ok := cause.Attacker.(*player.Player)
		if !ok {
			// TODO: support non-player attacker
			ctx.Cancel()
			return
		}

		a.mu.RLock()
		attackerPar, ok := a.u[user.Get(attacker).XUID()]
		a.mu.RUnlock()

		if !ok {
			ctx.Cancel()
			return
		}

		if death {
			a.BroadcastMessaget("killed.message.format", ctx.Val().Name(), attacker.Name())
			handledDeathMessage = true

			attackerPar.kills.Add(1)
			attackerPar.killStreak.Add(1)
		} else {
			par.StoreLastAttackedBy(user.Get(attacker).XUID())
		}
	case *entity.ProjectileDamageSource:
		owner, ok := cause.Owner.(*player.Player)
		if !ok {
			ctx.Cancel()
			return
		}

		a.mu.RLock()
		attackerPar, ok := a.u[user.Get(owner).XUID()]
		a.mu.RUnlock()

		if !ok {
			ctx.Cancel()
			return
		}

		if death {
			a.BroadcastMessaget("killed.shot.message.format", ctx.Val().Name(), owner.Name())
			handledDeathMessage = true

			attackerPar.kills.Add(1)
			attackerPar.killStreak.Add(1)
		} else {
			par.StoreLastAttackedBy(user.Get(owner).XUID())
		}
	case *entity.VoidDamageSource:
		if death {
			if lastAttackedBy := par.LastAttackedBy(); lastAttackedBy != "" {
				attacker, ok := a.u[lastAttackedBy]
				if ok {
					attacker.kills.Add(1)
					attacker.killStreak.Add(1)
					a.BroadcastMessaget("killed.void.message.format", ctx.Val().Name(), attacker.u.Name())
					break
				}
			}
			a.BroadcastMessaget("death.void.message.format", ctx.Val().Name())
			handledDeathMessage = true
		}
	}

	if !handledDeathMessage && death {
		a.BroadcastMessaget("death.message.format", ctx.Val().Name())
	}

	if death {
		par.deaths.Add(1)
		par.killStreak.Store(0)

		_ = a.Respawn(ctx.Val(), ctx.Val().Tx())
	} else {
		par.lastAttackedAt = time.Now()
		par.combatTimer.Store(11)
		helper.UpdatePlayerNameTagWithHealth(ctx.Val(), *damage)
	}
}

func (a *Arena) Icon() string {
	return a.icon
}

func (a *Arena) HandleBlockPlace(ctx *player.Context, pos cube.Pos, b world.Block) {
	if !a.allowBuild {
		ctx.Cancel()
		return
	}

	a.addPlacedBlock(pos, b)
}

func (a *Arena) HandleBlockBreak(ctx *player.Context, pos cube.Pos, drops *[]item.Stack, xp *int) {
	if !a.allowBuild {
		ctx.Cancel()
		return
	}

	a.removePlacedBlock(pos)
}

func (a *Arena) addPlacedBlock(pos cube.Pos, b world.Block) {
	a.placedBlocksMu.Lock()
	a.placedBlocks[pos] = b
	a.placedBlocksMu.Unlock()
}

func (a *Arena) removePlacedBlock(pos cube.Pos) {
	a.placedBlocksMu.Lock()
	delete(a.placedBlocks, pos)
	a.placedBlocksMu.Unlock()
}

func (a *Arena) HandleHeal(ctx *player.Context, health *float64, src world.HealingSource) {
	if ctx.Cancelled() {
		return
	}
	helper.UpdatePlayerNameTagWithHealth(ctx.Val(), *health)
}
