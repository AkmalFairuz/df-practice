package ffa

import (
	"errors"
	"github.com/akmalfairuz/df-practice/practice/helper"
	"github.com/akmalfairuz/df-practice/practice/kit"
	"github.com/akmalfairuz/df-practice/practice/kit/customitem"
	"github.com/akmalfairuz/df-practice/practice/user"
	"github.com/akmalfairuz/df-practice/translations"
	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/entity"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/particle"
	"github.com/df-mc/dragonfly/server/world/sound"
	"github.com/go-gl/mathgl/mgl64"
	"math/rand"
	"sync"
	"time"
)

type Arena struct {
	w *world.World

	u  map[string]*Participant
	mu sync.RWMutex

	spawns []helper.Location
	voidY  int

	dropAllowed bool
	k           kit.Kit

	icon string

	allowBuild bool

	placedBlocksMu sync.RWMutex
	placedBlocks   map[cube.Pos]placedBlockInfo

	zeroDamageExceptVoid bool

	attackCooldownTick int

	disableHPNameTag bool
	disableHunger    bool

	pearlCooldown time.Duration
}

type placedBlockInfo struct {
	placedAt    time.Time
	originBlock world.Block
	isBreaking  bool
}

func New(w *world.World) *Arena {
	return &Arena{
		w:             w,
		u:             make(map[string]*Participant),
		pearlCooldown: 15 * time.Second,
		placedBlocks:  make(map[cube.Pos]placedBlockInfo),
	}
}

func (a *Arena) applyConfig(config config) {
	a.spawns = config.SpawnLocations()
	a.voidY = config.VoidY

	if config.AttackCooldown != 0 {
		a.attackCooldownTick = config.AttackCooldown
	}

	if config.Icon != "" {
		a.icon = config.Icon
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

	a.w.Exec(func(tx *world.Tx) {
		for ent := range tx.Entities() {
			_ = ent.Close()
		}
	})

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
	a.w.Exec(func(tx *world.Tx) {
		for _, par := range a.Participants() {
			if currentTick%20 == 0 {
				par.combatTimer.Store(max(0, par.combatTimer.Load()-1))
				a.sendUserScoreboard(par, tx)
			}

			if p, ok := par.Player(tx); ok {
				if p.Immobile() && time.Since(par.lastSpawn.Load()) > time.Second {
					p.SetMobile()
				}

				if currentTick%2 == 0 {
					helper.UpdateXPBarCooldownDisplay(p, par.lastPearlThrow.Load(), a.pearlCooldown)
				}
			}
		}

		if currentTick%5 == 0 {
			if a.allowBuild {
				a.placedBlocksMu.RLock()
				placedBlocks := a.placedBlocks
				a.placedBlocksMu.RUnlock()

				for pos, info := range placedBlocks {
					diff := time.Since(info.placedAt)
					if diff > 15*time.Second {
						a.removePlacedBlock(pos)

						tx.AddParticle(pos.Vec3().Add(mgl64.Vec3{0.5, 0.5, 0.5}), particle.BlockBreak{Block: tx.Block(pos)})
						tx.PlaySound(pos.Vec3(), sound.BlockBreaking{Block: tx.Block(pos)})

						tx.SetBlock(pos, info.originBlock, &world.SetOpts{DisableBlockUpdates: true})
						continue
					}
					if diff > 8*time.Second {
						viewers := tx.Viewers(pos.Vec3())
						for _, viewer := range viewers {
							if info.isBreaking {
								viewer.ViewBlockAction(pos, block.ContinueCrackAction{
									BreakTime: time.Second * 7,
								})
							} else {
								viewer.ViewBlockAction(pos, block.StartCrackAction{
									BreakTime: time.Second * 7,
								})
							}
						}
						if !info.isBreaking {
							a.setPlacedBlockBreaking(pos, true)
						}
					}
				}
			}
		}
	})
}

func (a *Arena) sendUserScoreboard(p *Participant, tx *world.Tx) {
	combatInfo := "<empty>"
	if p.combatTimer.Load() > 0 {
		combatInfo = p.u.Translatef(translations.FfaScoreboardCombatTimer, p.combatTimer.Load())
	}

	p.u.SendScoreboard([]string{
		p.u.Translatef(translations.FfaScoreboardYourKills, p.kills.Load()),
		p.u.Translatef(translations.FfaScoreboardYourDeaths, p.deaths.Load()),
		p.u.Translatef(translations.FfaScoreboardYourStreak, p.killStreak.Load()),
		"",
		p.u.Translatef(translations.FfaScoreboardPlayers, len(a.u)),
		p.u.Translatef(translations.ScoreboardYourPing, p.u.Session().Latency().Milliseconds()),
		combatInfo,
	})
}

func (a *Arena) Join(p *player.Player, tx *world.Tx) error {
	u := user.Get(p)

	if u.CurrentGame() != nil {
		return errors.New("user already in game")
	}

	if u.CurrentFFAArena() != nil {
		return errors.New("user already in arena")
	}

	a.mu.RLock()
	if _, ok := a.u[u.XUID()]; ok {
		a.mu.RUnlock()
		return errors.New("user already in this arena, this should not happen after u.CurrentFFAArena() check before")
	}
	par := &Participant{u: u}
	par.lastSpawn.Store(time.Now())
	a.u[u.XUID()] = par
	a.mu.RUnlock()

	u.SetCurrentFFAArena(a)

	selectedSpawn := a.spawns[rand.Intn(len(a.spawns))]

	p.SetImmobile()
	tx.RemoveEntity(p)

	<-a.w.Exec(func(tx2 *world.Tx) {
		newP := tx2.AddEntity(p.H()).(*player.Player)
		selectedSpawn.TeleportPlayer(newP)

		helper.ResetPlayer(newP)
		_ = a.sendKit(newP)
		newP.SetGameMode(world.GameModeSurvival)

		if !a.disableHPNameTag {
			helper.UpdatePlayerNameTagWithHealth(newP, 0)
		}
	})

	a.sendUserScoreboard(par, tx)

	return nil
}

func (a *Arena) Respawn(p *player.Player) error {
	par, ok := a.ParticipantByXUID(user.Get(p).XUID())

	if !ok {
		return errors.New("user is not in this arena")
	}

	helper.ResetPlayer(p)
	selectedSpawn := a.spawns[rand.Intn(len(a.spawns))]
	selectedSpawn.TeleportPlayer(p)
	par.combatTimer.Store(0)
	par.lastPearlThrow.Store(time.Unix(0, 0))
	par.lastAttackedAt.Store(time.Unix(0, 0))
	par.lastSpawn.Store(time.Now())
	if !a.disableHPNameTag {
		helper.UpdatePlayerNameTagWithHealth(p, 0)
	}
	p.SetImmobile()
	_ = a.sendKit(p)
	return nil
}

func (a *Arena) sendKit(p *player.Player) error {
	if a.k != nil {
		kit.Apply(a.k, p)
	}
	return nil
}

// Quit removes a player from the arena. Caller should teleport the player to lobby.
func (a *Arena) Quit(p *player.Player) error {
	p.SetMobile()
	u := user.Get(p)

	a.mu.RLock()
	_, ok := a.u[u.XUID()]
	a.mu.RUnlock()
	if !ok {
		return errors.New("user is not in this arena")
	}

	a.mu.Lock()
	delete(a.u, u.XUID())
	a.mu.Unlock()

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
	if helper.InvalidPlayerCtxWorld(ctx, a.w) {
		return
	}

	if immune {
		return
	}

	par, ok := a.u[ctx.Val().XUID()]
	if !ok {
		return
	}

	if a.zeroDamageExceptVoid {
		if _, ok := src.(entity.VoidDamageSource); !ok {
			*damage = 0
		}
	}

	if time.Since(par.lastSpawn.Load()) < 3*time.Second {
		ctx.Cancel()
		if _, ok := src.(entity.AttackDamageSource); ok {
			ctx.Val().Tx().AddParticle(ctx.Val().Position(), particle.Lava{})
		}
		return
	}

	death := *damage >= ctx.Val().Health()

	if death {
		ctx.Cancel()
	}

	handledDeathMessage := false

	switch cause := src.(type) {
	case entity.AttackDamageSource:
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

		if a.attackCooldownTick != 0 {
			*immunity = time.Millisecond * 50 * time.Duration(a.attackCooldownTick)
		}

		if death {
			a.BroadcastMessaget(translations.KilledMessageFormat, ctx.Val().Name(), attacker.Name())
			handledDeathMessage = true

			a.OnKill(attacker, attackerPar)
		} else {
			par.StoreLastAttackedBy(user.Get(attacker).XUID())
		}
	case entity.ProjectileDamageSource:
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
			if ctx.Val().H() == owner.H() {
				a.BroadcastMessaget(translations.KilledSelfShotMessageFormat, ctx.Val().Name())
				handledDeathMessage = true
				break
			}
			a.BroadcastMessaget(translations.KilledShotMessageFormat, ctx.Val().Name(), owner.Name())
			handledDeathMessage = true

			a.OnKill(owner, attackerPar)
		} else {
			par.StoreLastAttackedBy(user.Get(owner).XUID())
		}
	case entity.VoidDamageSource:
		if death {
			if lastAttackedBy := par.LastAttackedBy(); lastAttackedBy != "" {
				attacker, ok := a.u[lastAttackedBy]
				if ok {
					if ctx.Val().Name() == attacker.u.Name() {
						a.BroadcastMessaget(translations.KilledSelfVoidMessageFormat, ctx.Val().Name())
						handledDeathMessage = true
						break
					}
					// Should using EntityHandle later

					attacker.kills.Add(1)
					attacker.killStreak.Add(1)
					a.BroadcastMessaget(translations.KilledVoidMessageFormat, ctx.Val().Name(), attacker.u.Name())
					handledDeathMessage = true
					break
				}
			}
			a.BroadcastMessaget(translations.DeathVoidMessageFormat, ctx.Val().Name())
			handledDeathMessage = true
		}
	}

	if !handledDeathMessage && death {
		a.BroadcastMessaget(translations.DeathMessageFormat, ctx.Val().Name())
	}

	if death {
		par.deaths.Add(1)
		par.killStreak.Store(0)

		helper.LogErrors(a.Respawn(ctx.Val()))
	} else {
		par.lastAttackedAt.Store(time.Now())
		par.combatTimer.Store(16)
		if !a.disableHPNameTag {
			helper.UpdatePlayerNameTagWithHealth(ctx.Val(), 0-*damage)
		}
	}
}

func (a *Arena) Icon() string {
	return a.icon
}

func (a *Arena) HandleBlockPlace(ctx *player.Context, pos cube.Pos, b world.Block) {
	if helper.InvalidPlayerCtxWorld(ctx, a.w) {
		return
	}

	if !a.allowBuild {
		ctx.Cancel()
		return
	}
	if a.hasPlacedBlock(pos) {
		ctx.Cancel()
		user.Messaget(ctx.Val(), translations.ErrorPlaceBlock)
		return
	}
	for _, spawn := range a.spawns {
		if spawn.ToMgl64Vec3().Sub(pos.Vec3()).Len() < 2 {
			ctx.Cancel()
			user.Messaget(ctx.Val(), translations.ErrorPlaceBlock)
			return
		}
	}
	currentBlock := ctx.Val().Tx().Block(pos)
	if _, ok := currentBlock.(block.Air); !ok {
		ctx.Cancel()
		user.Messaget(ctx.Val(), translations.ErrorPlaceBlock)
		return
	}

	a.addPlacedBlock(pos, currentBlock)
}

func (a *Arena) HandleBlockBreak(ctx *player.Context, pos cube.Pos, drops *[]item.Stack, xp *int) {
	if helper.InvalidPlayerCtxWorld(ctx, a.w) {
		return
	}

	if !a.allowBuild {
		ctx.Cancel()
		return
	}
	if !a.hasPlacedBlock(pos) {
		ctx.Cancel()
		return
	}

	a.removePlacedBlock(pos)
}

func (a *Arena) hasPlacedBlock(pos cube.Pos) bool {
	a.placedBlocksMu.RLock()
	_, ok := a.placedBlocks[pos]
	a.placedBlocksMu.RUnlock()
	return ok
}

func (a *Arena) addPlacedBlock(pos cube.Pos, b world.Block) {
	a.placedBlocksMu.Lock()
	a.placedBlocks[pos] = placedBlockInfo{placedAt: time.Now(), originBlock: b}
	a.placedBlocksMu.Unlock()
}

func (a *Arena) removePlacedBlock(pos cube.Pos) {
	a.placedBlocksMu.Lock()
	delete(a.placedBlocks, pos)
	a.placedBlocksMu.Unlock()
}

func (a *Arena) setPlacedBlockBreaking(pos cube.Pos, breaking bool) {
	a.placedBlocksMu.Lock()
	info, ok := a.placedBlocks[pos]
	if ok {
		info.isBreaking = breaking
		a.placedBlocks[pos] = info
	}
	a.placedBlocksMu.Unlock()
}

func (a *Arena) HandleHeal(ctx *player.Context, health *float64, src world.HealingSource) {
	if helper.InvalidPlayerCtxWorld(ctx, a.w) {
		return
	}

	if ctx.Cancelled() {
		return
	}

	if !a.disableHPNameTag {
		helper.UpdatePlayerNameTagWithHealth(ctx.Val(), *health)
	}
}

func (a *Arena) HandleFoodLoss(ctx *player.Context, from int, to *int) {
	if helper.InvalidPlayerCtxWorld(ctx, a.w) {
		return
	}

	if a.disableHunger {
		ctx.Cancel()
	}
}

func (a *Arena) HandleItemDrop(ctx *player.Context, s item.Stack) {
	if helper.InvalidPlayerCtxWorld(ctx, a.w) {
		return
	}

	if !a.dropAllowed {
		ctx.Cancel()
	}
}

func (a *Arena) HandleMove(ctx *player.Context, pos mgl64.Vec3, rot cube.Rotation) {
	if helper.InvalidPlayerCtxWorld(ctx, a.w) {
		return
	}

	if pos.Y() < float64(a.voidY) {
		ctx.Val().Hurt(1000, entity.VoidDamageSource{})
	}
}

func (a *Arena) HandleItemUse(ctx *player.Context) {
	if helper.InvalidPlayerCtxWorld(ctx, a.w) {
		return
	}

	par, ok := a.ParticipantByXUID(ctx.Val().XUID())
	if !ok {
		return
	}

	mainHand, _ := ctx.Val().HeldItems()

	if mainHand.Comparable(item.NewStack(customitem.NoDamageEnderPearl{}, 1)) && !ctx.Cancelled() {
		lastPearlThrow := par.lastPearlThrow.Load()
		if time.Since(lastPearlThrow) < a.pearlCooldown {
			par.u.Messaget(translations.ErrorCooldownPearl, time.Until(lastPearlThrow.Add(a.pearlCooldown)).Seconds())
			ctx.Cancel()
			return
		}
		par.lastPearlThrow.Store(time.Now())
	}
}

func (a *Arena) OnKill(p *player.Player, par *Participant) {
	par.OnKill()
	p.PlaySound(sound.Experience{})
	helper.ClearAllPlayerInv(p)
	_ = a.sendKit(p)
}

func (a *Arena) HandleStartBreak(ctx *player.Context, pos cube.Pos) {
	if helper.InvalidPlayerCtxWorld(ctx, a.w) {
		return
	}

	if !a.allowBuild {
		ctx.Cancel()
		return
	}

	a.placedBlocksMu.RLock()
	info, ok := a.placedBlocks[pos]
	a.placedBlocksMu.RUnlock()

	if !ok {
		ctx.Cancel()
		return
	}

	if info.isBreaking {
		a.placedBlocksMu.Lock()
		a.placedBlocks[pos] = placedBlockInfo{placedAt: time.Now(), originBlock: info.originBlock}
		a.placedBlocksMu.Unlock()

		for _, viewer := range ctx.Val().Tx().Viewers(pos.Vec3()) {
			viewer.ViewBlockAction(pos, block.StopCrackAction{})
		}

		ctx.Cancel()
	}
}
