package game

import (
	"errors"
	"github.com/akmalfairuz/df-practice/internal/meta"
	"github.com/akmalfairuz/df-practice/practice/game/igame"
	"github.com/akmalfairuz/df-practice/practice/helper"
	"github.com/akmalfairuz/df-practice/practice/lang"
	"github.com/akmalfairuz/df-practice/practice/user"
	"github.com/df-mc/atomic"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/item/inventory"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/player/title"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
	"log/slog"
	"math/rand"
	"os"
	"sync"
	"time"
)

func quitItem(p *player.Player) item.Stack {
	return item.NewStack(item.DragonBreath{}, 1).WithCustomName(lang.Translatef(user.Lang(p), "game.item.quit.name")).WithValue("game_item", "quit")
}

func init() {
	if err := helper.RemoveDir(gameWorldsPath); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			panic(err)
		}
	}
	if err := os.Mkdir(gameWorldsPath, 0755); err != nil {
		panic(err)
	}
}

type Game struct {
	log *slog.Logger

	id string

	pMu sync.RWMutex
	p   map[string]igame.IParticipant

	currentTick atomic.Uint64

	state atomic.Value[State]

	impl  igame.Impl
	pImpl igame.IParticipant

	w    *world.World
	wDir string

	closed atomic.Bool
	once   sync.Once

	closeHook func()

	tickQueue chan struct{}
	closing   chan struct{}
}

const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func generateID() string {
	b := make([]byte, 8)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func (g *Game) Load() error {
	g.closed.Store(false)
	g.tickQueue = make(chan struct{})
	g.closing = make(chan struct{})

	g.impl.OnInit()

	go g.handleTick()
	go g.startTicking()
	return nil
}

func (g *Game) startTicking() {
	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if g.closed.Load() {
				g.closing <- struct{}{}
				return
			}
			g.tickQueue <- struct{}{}
		}
	}
}

func (g *Game) handleTick() {
	for {
		select {
		case <-g.tickQueue:
			g.OnTick()
		case <-g.closing:
			return
		}
	}
}

func (g *Game) OnTick() {
	g.impl.OnTick()
	g.currentTick.Add(1)

	switch g.state.Load() {
	case StateWaiting:
		if len(g.p) < g.impl.MinimumParticipants() {
			g.currentTick.Store(0)
		}

		if g.currentTick.Load()%20 == 0 {
			for _, p := range g.Participants() {
				u := p.User()

				lines := make([]string, 0)
				lines = append(lines, u.Translatef("game.waiting.scoreboard.players", len(g.p), g.impl.MaxParticipants()))
				lines = append(lines, "")
				if len(g.p) >= g.impl.MinimumParticipants() {
					lines = append(lines, u.Translatef("game.waiting.scoreboard.starting.in", g.impl.WaitingTime()-int(g.currentTick.Load()/20)))
				} else {
					lines = append(lines, u.Translatef("game.waiting.scoreboard.waiting.for.players"))
				}

				u.SendScoreboard(lines)
			}

			remaining := g.impl.WaitingTime() - int(g.currentTick.Load()/20)
			if remaining <= 5 && remaining > 0 {
				g.Messaget("game.waiting.starting.message", remaining)
				g.Sound("random.click")
			}
		}

		if g.CurrentTick() >= uint64(g.impl.WaitingTime())*20 {
			g.Start()
		}
	case StatePlaying:
		if g.CurrentTick() >= uint64(g.impl.PlayingTime())*20 {
			g.End()
			return
		}
	case StateEnding:
		if g.currentTick.Load()%20 == 0 {
			for _, p := range g.Participants() {
				u := p.User()

				lines := make([]string, 0)
				lines = append(lines, u.Translatef("game.ending.scoreboard.stopping.in", g.impl.EndingTime()-int(g.currentTick.Load()/20)))

				u.SendScoreboard(lines)
			}
		}
		if g.currentTick.Load() >= uint64(g.impl.EndingTime())*20 {
			g.Stop()
			return
		}
	}
}

func (g *Game) Start() {
	g.state.Store(StatePlaying)
	g.currentTick.Store(0)
	g.Sound("mob.blaze.shoot")
	g.Messaget("game.started.message")
	g.w.Exec(func(tx *world.Tx) {
		for ent := range tx.Entities() {
			p, ok := ent.(*player.Player)
			if !ok {
				continue
			}

			_, ok = g.ParticipantByXUID(p.XUID())
			if !ok {
				continue
			}

			helper.ResetPlayer(p)
		}
	})
	g.impl.OnStart()
}

func (g *Game) End() {
	if g.IsEnding() {
		return
	}
	g.currentTick.Store(0)
	g.state.Store(StateEnding)
	g.Sound("random.explode")
	g.Messaget("game.ended.message")
	g.w.Exec(func(tx *world.Tx) {
		for ent := range tx.Entities() {
			p, ok := ent.(*player.Player)
			if !ok {
				continue
			}

			par, ok := g.ParticipantByXUID(p.XUID())
			if !ok {
				continue
			}

			helper.ResetPlayer(p)
			if par.IsPlaying() {
				p.SetGameMode(world.GameModeAdventure)
			}
			_ = p.Inventory().SetItem(8, quitItem(p))
			_ = p.SetHeldSlot(0)
		}

	})

	g.impl.OnEnd()
}

func (g *Game) Stop() {
	g.once.Do(func() {
		defer g.closed.Store(true)

		g.impl.OnStop()

		<-g.w.Exec(func(tx *world.Tx) {
			for ent := range tx.Entities() {
				p, ok := ent.(*player.Player)
				if !ok {
					continue
				}

				if _, ok := g.p[p.XUID()]; ok {
					helper.LogErrors(g.Quit(p))
				}
			}
		})

		helper.LogErrors(g.w.Close())
		helper.LogErrors(helper.RemoveDir(g.wDir))

		if g.closeHook != nil {
			(g.closeHook)()
		}
	})
}

func (g *Game) ID() string {
	return g.id
}

func (g *Game) Participants() map[string]igame.IParticipant {
	g.pMu.RLock()
	defer g.pMu.RUnlock()
	return g.p
}

func (g *Game) PlayingParticipants() map[string]igame.IParticipant {
	g.pMu.RLock()
	defer g.pMu.RUnlock()

	par := make(map[string]igame.IParticipant)
	for xuid, p := range g.p {
		if !p.IsSpectating() {
			par[xuid] = p
		}
	}
	return par
}

func (g *Game) SpectatingParticipants() map[string]igame.IParticipant {
	g.pMu.RLock()
	defer g.pMu.RUnlock()

	par := make(map[string]igame.IParticipant)
	for xuid, p := range g.p {
		if p.IsSpectating() {
			par[xuid] = p
		}
	}
	return par
}

func (g *Game) ParticipantByXUID(xuid string) (igame.IParticipant, bool) {
	g.pMu.RLock()
	defer g.pMu.RUnlock()
	p, ok := g.p[xuid]
	return p, ok
}

func (g *Game) Join(p *player.Player) error {
	if g.state.Load() != StateWaiting {
		return errors.New("game already started")
	}

	g.pMu.Lock()

	if len(g.p) >= g.impl.MaxParticipants() {
		g.pMu.Unlock()
		return errors.New("max participants reached")
	}

	if _, ok := g.p[p.XUID()]; ok {
		g.pMu.Unlock()
		return errors.New("already joined")
	}

	if err := g.impl.OnJoin(p); err != nil {
		g.pMu.Unlock()
		return err
	}

	par := g.newParticipant(p)
	g.p[p.XUID()] = par
	g.pMu.Unlock()

	user.Get(p).SetCurrentGame(g)

	p.Tx().RemoveEntity(p)
	<-g.w.Exec(func(tx *world.Tx) {
		newP := tx.AddEntity(p.H()).(*player.Player)

		helper.ResetPlayer(newP)
		_ = p.Inventory().SetItem(8, quitItem(newP))

		p.SetGameMode(world.GameModeAdventure)
		g.impl.OnJoined(par, newP)

		g.Messaget("game.waiting.join.message", user.Get(newP).Name(), len(g.p), g.impl.MaxParticipants())
	})
	return nil
}

func (g *Game) Quit(p *player.Player) error {
	g.pMu.RLock()
	if _, ok := g.p[p.XUID()]; !ok {
		g.pMu.RUnlock()
		return errors.New("not joined")
	}
	g.pMu.RUnlock()

	if g.IsWaiting() {
		g.Messaget("game.waiting.left.message", user.Get(p).Name(), len(g.p)-1, g.impl.MaxParticipants())
	}
	g.impl.OnQuit(p)
	helper.ResetPlayer(p)
	p.SetGameMode(world.GameModeAdventure)
	p.SetMobile()

	g.pMu.Lock()
	delete(g.p, p.XUID())
	g.pMu.Unlock()

	if g.IsPlaying() {
		g.impl.CheckEnd()
	}

	user.Get(p).SetCurrentGame(nil)

	//p.Tx().RemoveEntity(p)
	meta.Get("lobby").(interface {
		Spawn(p *player.Player)
	}).Spawn(p)
	return nil
}

func (g *Game) State() State {
	return g.state.Load()
}

func (g *Game) IsPlaying() bool {
	return g.state.Load() == StatePlaying
}

func (g *Game) IsWaiting() bool {
	return g.state.Load() == StateWaiting
}

func (g *Game) IsEnding() bool {
	return g.state.Load() == StateEnding
}

func (g *Game) newParticipant(p *player.Player) *Participant {
	return &Participant{
		xuid: p.XUID(),
	}
}

func (g *Game) HandleItemUse(ctx *player.Context) {
	par, _ := g.ParticipantByXUID(ctx.Val().XUID())

	if !g.IsPlaying() || par.IsSpectating() {
		ctx.Cancel()

		mainHand, _ := ctx.Val().HeldItems()
		v, ok := mainHand.Value("game_item")
		if ok {
			switch v {
			case "quit":
				helper.LogErrors(g.Quit(ctx.Val()))
				return
			}
		}
	}

	g.impl.HandleItemUse(ctx)
}

func (g *Game) Messaget(translationName string, args ...any) {
	for _, p := range g.Participants() {
		u := user.GetByXUID(p.XUID())
		if u != nil {
			u.Messaget(translationName, args...)
		}
	}
}

func (g *Game) Sound(name string) {
	g.w.Exec(func(tx *world.Tx) {
		for _, p := range g.Participants() {
			ent, ok := p.User().EntityHandle().Entity(tx)
			if !ok {
				continue
			}
			playSound := &packet.PlaySound{
				SoundName: name,
				Position:  helper.Mgl64Vec3ToMgl32Vec3(ent.Position()),
				Volume:    0.5,
				Pitch:     1,
			}
			_ = p.User().Conn().WritePacket(playSound)
		}
	})
}

func (g *Game) World() *world.World {
	return g.w
}

func (g *Game) HandleAttackEntity(ctx *player.Context, e world.Entity, force *float64, height *float64, critical *bool) {
	if !g.IsPlaying() {
		ctx.Cancel()
		return
	}
	if par, ok := g.ParticipantByXUID(ctx.Val().XUID()); ok && par.IsSpectating() {
		ctx.Cancel()
		return
	}

	g.impl.HandleAttackEntity(ctx, e, force, height, critical)
}

func (g *Game) HandleHurt(ctx *player.Context, damage *float64, immune bool, immunity *time.Duration, src world.DamageSource) {
	if !g.IsPlaying() {
		ctx.Cancel()
		return
	}
	if par, ok := g.ParticipantByXUID(ctx.Val().XUID()); ok && par.IsSpectating() {
		ctx.Cancel()
		return
	}

	g.impl.HandleHurt(ctx, damage, immune, immunity, src)
}

func (g *Game) HandleHeal(ctx *player.Context, health *float64, src world.HealingSource) {
	g.impl.HandleHeal(ctx, health, src)
}

func (g *Game) HandleFoodLoss(ctx *player.Context, from int, to *int) {
	if !g.IsPlaying() {
		ctx.Cancel()
		return
	}

	g.impl.HandleFoodLoss(ctx, from, to)
}

func (g *Game) HandleBlockBreak(ctx *player.Context, pos cube.Pos, drops *[]item.Stack, xp *int) {
	if !g.IsPlaying() {
		ctx.Cancel()
		return
	}

	g.impl.HandleBlockBreak(ctx, pos, drops, xp)
}

func (g *Game) HandleBlockPlace(ctx *player.Context, pos cube.Pos, b world.Block) {
	if !g.IsPlaying() {
		ctx.Cancel()
		return
	}

	g.impl.HandleBlockPlace(ctx, pos, b)
}

func (g *Game) HandleMove(ctx *player.Context, pos mgl64.Vec3, rot cube.Rotation) {
	g.impl.HandleMove(ctx, pos, rot)
}

func (g *Game) HandleItemUseOnEntity(ctx *player.Context, e world.Entity) {
	g.impl.HandleItemUseOnEntity(ctx, e)
}

func (g *Game) HandleDrop(ctx *inventory.Context, slot int, stack item.Stack) {
	if !g.IsPlaying() {
		ctx.Cancel()
		return
	}

	g.impl.HandleDrop(ctx, slot, stack)
}

func (g *Game) HandlePlace(ctx *inventory.Context, slot int, stack item.Stack) {
	if !g.IsPlaying() {
		ctx.Cancel()
		return
	}

	g.impl.HandlePlace(ctx, slot, stack)
}

func (g *Game) HandleTake(ctx *inventory.Context, slot int, stack item.Stack) {
	if !g.IsPlaying() {
		ctx.Cancel()
		return
	}

	g.impl.HandleTake(ctx, slot, stack)
}

func (g *Game) SetSpectator(p *player.Player) {
	par, ok := g.ParticipantByXUID(p.XUID())
	if !ok {
		return
	}

	par.(*Participant).state = ParticipantStateSpectating
	helper.ResetPlayer(p)
	p.SetMobile()
	p.SetGameMode(world.GameModeSpectator)
	_ = p.Inventory().SetItem(8, quitItem(p))
	_ = p.SetHeldSlot(0)

	p.SendTitle(title.New(lang.Translatef(user.Lang(p), "game.you.died.title")))

	g.impl.CheckEnd()
}

func (g *Game) CurrentTick() uint64 {
	return g.currentTick.Load()
}

func (g *Game) Players(tx *world.Tx) []*player.Player {
	players := make([]*player.Player, 0)
	for ent := range tx.Entities() {
		p, ok := ent.(*player.Player)
		if ok {
			players = append(players, p)
		}
	}
	return players
}

func (g *Game) SetCloseHook(hook func()) {
	g.closeHook = hook
}

func (g *Game) HandleItemPickup(ctx *player.Context, i *item.Stack) {
	if !g.IsPlaying() {
		ctx.Cancel()
		return
	}

	par, ok := g.ParticipantByXUID(ctx.Val().XUID())
	if ok {
		if par.IsSpectating() {
			ctx.Cancel()
			return
		}
	}
}
