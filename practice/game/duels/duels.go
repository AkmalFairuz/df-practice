package duels

import (
	"github.com/akmalfairuz/df-practice/practice/game/gamedefaults"
	"github.com/akmalfairuz/df-practice/practice/game/igame"
	"github.com/akmalfairuz/df-practice/practice/helper"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"sync"
	"time"
)

type spawnInfo struct {
	usedBy string
	loc    helper.Location
}

type Duels struct {
	igame.Nop

	g igame.IGame

	spawns   []*spawnInfo
	spawnsMu sync.Mutex

	onSendKit func(p *player.Player) error
}

func (d *Duels) Create(g igame.IGame) {
	d.spawns = []*spawnInfo{
		{loc: helper.Location{X: 0 + 0.5, Y: 68.2, Z: 14 + 0.5, Yaw: -180, Pitch: 0}},
		{loc: helper.Location{X: 0 + 0.5, Y: 68.2, Z: -14 + 0.5, Yaw: 0, Pitch: 0}},
	}

	d.g = g
}

func (d *Duels) PlayingTime() int {
	return 900
}

func (d *Duels) WaitingTime() int {
	return 8
}

func (d *Duels) EndingTime() int {
	return 10
}

func (d *Duels) MaxParticipants() int {
	return 2
}

func (d *Duels) MinimumParticipants() int {
	return 2
}

func (d *Duels) Name() string {
	return "Duels"
}

func (d *Duels) OnJoin(p *player.Player) error {
	return nil
}

func (d *Duels) OnJoined(par igame.IParticipant, p *player.Player) {
	d.spawnsMu.Lock()
	defer d.spawnsMu.Unlock()

	for i, spawn := range d.spawns {
		if spawn.usedBy == "" {
			d.spawns[i].usedBy = par.User().XUID()
			d.spawns[i].loc.TeleportPlayer(p)
			break
		}
	}

	p.SetImmobile()
}

func (d *Duels) OnQuit(p *player.Player) {
	if d.g.IsWaiting() {
		d.spawnsMu.Lock()
		defer d.spawnsMu.Unlock()

		for i, spawn := range d.spawns {
			if spawn.usedBy == p.XUID() {
				d.spawns[i].usedBy = ""
				break
			}
		}
	}
}

func (d *Duels) OnInit() {

}

func (d *Duels) OnStart() {
	<-d.g.World().Exec(func(tx *world.Tx) {
		for _, p := range d.g.Players(tx) {
			p.SetMobile()
			p.SetGameMode(world.GameModeSurvival)

			helper.LogErrors((d.onSendKit)(p))
		}
	})
}

func (d *Duels) OnEnd() {

}

func (d *Duels) OnStop() {

}

func (d *Duels) OnTick() {
	if d.g.IsPlaying() {
		if d.g.CurrentTick()%20 == 0 {
			for _, p := range d.g.Participants() {
				opponent, ok := d.opponent(p)
				opponentPing := 0
				if ok {
					opponentPing = opponent.User().Ping()
				}

				p.User().SendScoreboard([]string{
					p.User().Translatef("game.scoreboard.time.left", helper.FormatTime(int(d.g.CurrentTick()/20))),
					"",
					p.User().Translatef("scoreboard.their.ping", opponentPing),
					p.User().Translatef("scoreboard.your.ping", p.User().Ping()),
				})
			}
		}
	}
}

func (d *Duels) CheckEnd() {
	if len(d.g.PlayingParticipants()) <= 1 {
		d.g.End()
	}
}

func (d *Duels) HandleHurt(ctx *player.Context, damage *float64, immune bool, immunity *time.Duration, src world.DamageSource) {
	death := *damage >= ctx.Val().Health()
	if death {
		ctx.Cancel()

		gamedefaults.HandleKillMessage(d.g, ctx.Val(), src)
	}
	if death {
		d.g.SetSpectator(ctx.Val())
	}
}

func (d *Duels) opponent(p igame.IParticipant) (igame.IParticipant, bool) {
	for _, par := range d.g.Participants() {
		if par != p {
			return par, true
		}
	}
	return nil, false
}

func (d *Duels) HandleBlockBreak(ctx *player.Context, pos cube.Pos, drops *[]item.Stack, xp *int) {
	ctx.Cancel()
}

func (d *Duels) HandleBlockPlace(ctx *player.Context, pos cube.Pos, b world.Block) {
	ctx.Cancel()
}

func (d *Duels) Game() igame.IGame {
	return d.g
}

func (d *Duels) SetKit(kit func(p *player.Player) error) {
	d.onSendKit = kit
}

// Compile-time check to ensure that Duels implements game.Impl.
var _ igame.Impl = &Duels{}
