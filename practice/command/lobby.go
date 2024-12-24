package command

import (
	"github.com/akmalfairuz/df-practice/practice/ffa"
	"github.com/akmalfairuz/df-practice/practice/helper"
	"github.com/akmalfairuz/df-practice/practice/lobby"
	"github.com/akmalfairuz/df-practice/practice/user"
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
)

type Lobby struct {
}

func (l Lobby) Run(s cmd.Source, o *cmd.Output, tx *world.Tx) {
	u := user.Get(s.(*player.Player))

	ffaArena := u.CurrentFFAArena()
	if ffaArena != nil {
		a := ffaArena.(*ffa.Arena)
		par, _ := a.ParticipantByXUID(u.XUID())
		if par.InCombat() {
			o.Error(u.Translatef("error.lobby.in.combat"))
			return
		}
		helper.LogErrors(a.Quit(s.(*player.Player)))
	}

	lobby.Instance().Spawn(s.(*player.Player))
}

func (l Lobby) Allow(s cmd.Source) bool {
	return isPlayer(s)
}
