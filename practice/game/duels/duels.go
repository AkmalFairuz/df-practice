package duels

import (
	"github.com/akmalfairuz/df-practice/practice/game"
	"github.com/akmalfairuz/df-practice/practice/game/nop"
	"github.com/akmalfairuz/df-practice/practice/helper"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"log/slog"
	"time"
)

type Duels struct {
	nop.Game

	g *game.Game
}

func New(log *slog.Logger, mapDir string) *Duels {
	impl := &Duels{}

	gConf := game.Config{
		Log:    log,
		MapDir: mapDir,
		Impl:   impl,
	}

	g, err := gConf.New()
	if err != nil {
		panic(err)
	}

	impl.g = g
	return impl
}

func (d *Duels) PlayingTime() int {
	return 900
}

func (d *Duels) WaitingTime() int {
	return 6
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

func (d *Duels) OnJoined(par *game.Participant, p *player.Player) {

}

func (d *Duels) OnQuit(p *player.Player) {

}

func (d *Duels) OnInit() {

}

func (d *Duels) OnStart() {

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
	}
	// TODO: Handle death message
	if death {
		d.g.SetSpectator(ctx.Val())
	}
}

func (d *Duels) opponent(p *game.Participant) (*game.Participant, bool) {
	for _, par := range d.g.Participants() {
		if par != p {
			return par, true
		}
	}
	return nil, false
}

// Compile-time check to ensure that Duels implements game.Impl.
var _ game.Impl = &Duels{}
