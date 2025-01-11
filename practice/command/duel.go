package command

import (
	"github.com/akmalfairuz/df-practice/practice/game/duelsmanager"
	"github.com/akmalfairuz/df-practice/practice/user"
	"github.com/akmalfairuz/df-practice/translations"
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"time"
)

type Duel struct {
	onlyPlayer

	Target onlineTarget `cmd:"target"`
}

func tryAcceptDuelReq(u *user.User, from *user.User, tx *world.Tx) bool {
	info := from.DuelRequestTo()
	if info.TargetXUID != u.XUID() {
		return false
	}
	if time.Since(info.RequestAt) > time.Second*60 {
		return false
	}
	if !from.InLobby() || from.CurrentGame() != nil || from.CurrentFFAArena() != nil {
		return false
	}
	if !u.InLobby() || u.CurrentGame() != nil || u.CurrentFFAArena() != nil {
		return false
	}
	from.SetDuelRequestTo("")
	u.SetDuelRequestTo("")

	g := duelsmanager.ClassicRaw()
	if err := g.Game().Load(); err != nil {
		panic(err)
	}
	g.Game().SetDuelRequestMode(true)

	fromEnt, _ := from.Player(tx)
	uEnt, _ := u.Player(tx)

	from.Messaget(translations.DuelRequestAccepted, u.Name())

	_ = g.Game().Join(fromEnt)
	_ = g.Game().Join(uEnt)

	return true
}

func (d Duel) Run(s cmd.Source, o *cmd.Output, tx *world.Tx) {
	u := user.Get(s.(*player.Player))
	if !u.InLobby() {
		u.Messaget(translations.DuelRequestMustInLobby)
		return
	}

	targetU := d.Target.User()
	if targetU == u {
		u.Messaget(translations.CommandAnErrorOccurred)
		return
	}
	if tryAcceptDuelReq(u, targetU, tx) {
		return
	}

	u.OnSendDuelRequest(targetU)
	targetU.OnReceiveDuelRequest(u)
}
