package lobby

import (
	"github.com/akmalfairuz/df-practice/practice/helper"
	"github.com/akmalfairuz/df-practice/practice/lang"
	"github.com/akmalfairuz/df-practice/practice/user"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"sync/atomic"
	"time"
)

type Lobby struct {
	w *world.World

	closed atomic.Bool
}

var instance *Lobby

func Instance() *Lobby {
	return instance
}

func New(w *world.World) *Lobby {
	ret := &Lobby{
		w: w,
	}
	instance = ret
	return ret
}

const (
	lobbyItemIndexKey = "lobby_item"
)

func (l *Lobby) Init() error {
	l.closed.Store(false)
	l.w.Handle(newLobbyWorldHandler(l))

	l.w.StopThundering()
	l.w.StopRaining()
	l.w.StopWeatherCycle()
	l.w.StopTime()
	l.w.SetDifficulty(world.DifficultyEasy)

	l.w.SetTime(3000)

	go l.startTicking()

	return nil
}

func (l *Lobby) startTicking() {
	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()

	currentTick := int64(0)

	for {
		select {
		case <-ticker.C:
			if l.closed.Load() {
				return
			}
			currentTick++
			l.handleTick(currentTick)
		}
	}
}

func (l *Lobby) handleTick(currentTick int64) {
	<-l.w.Exec(func(tx *world.Tx) {
		users := l.users(tx)
		for _, u := range users {
			l.onUserTick(u, tx, currentTick)
		}
	})
}

func (l *Lobby) onUserTick(u *user.User, tx *world.Tx, currentTick int64) {
	if currentTick%20 == 0 {
		l.sendUserScoreboard(u, tx)
	}
}

func (l *Lobby) sendUserScoreboard(u *user.User, tx *world.Tx) {
	u.SendScoreboard([]string{
		u.Translatef("lobby.scoreboard.players.online.count", user.Count()),
		u.Translatef("scoreboard.your.ping", u.Session().Latency().Milliseconds()),
	})
}

func (l *Lobby) World() *world.World {
	return l.w
}

func (l *Lobby) users(tx *world.Tx) []*user.User {
	users := make([]*user.User, 0)
	for p := range tx.Players() {
		u := user.Get(p.(*player.Player))
		if u != nil {
			users = append(users, u)
		}
	}
	return users
}

func (l *Lobby) sendLobbyItems(p *player.Player) {
	u := user.Get(p)

	_ = p.Inventory().SetItem(0, item.NewStack(item.Sword{Tier: item.ToolTierDiamond}, 1).WithCustomName(lang.Translatef(u.Lang(), "lobby.item.play.ffa.name")).WithValue(lobbyItemIndexKey, 0))
	//_ = p.Inventory().SetItem(1, item.NewStack(item.Sword{Tier: item.ToolTierIron}, 1).WithCustomName(lang.Translatef(u.Lang(), "lobby.item.play.duels.name")).WithValue(lobbyItemIndexKey, 1))
	//_ = p.Inventory().SetItem(8, item.NewStack(block.Skull{Type: block.PlayerHead()}, 1).WithCustomName(lang.Translatef(u.Lang(), "lobby.item.profile.name")).WithValue(lobbyItemIndexKey, 8))
}

func (l *Lobby) Spawn(p *player.Player) {
	if p.Tx().World() == l.w {
		p.Teleport(l.w.Spawn().Vec3())
		helper.ResetPlayer(p)
		l.sendLobbyItems(p)
		l.sendUserScoreboard(user.Get(p), p.Tx())
	} else {
		ent, _ := p.H().Entity(p.Tx())
		p.Tx().RemoveEntity(ent)

		l.w.Exec(func(tx *world.Tx) {
			newP := tx.AddEntity(p.H()).(*player.Player)
			helper.ResetPlayer(newP)
			l.sendLobbyItems(newP)
			l.sendUserScoreboard(user.Get(newP), tx)
			newP.SetGameMode(world.GameModeAdventure)
		})
	}
}

func (l *Lobby) IsInLobby(p *player.Player) bool {
	return p.Tx().World() == l.w
}

func (l *Lobby) HandleItemUse(ctx *player.Context) {
	if ctx.Val().Tx().World() != l.w || ctx.Cancelled() {
		return
	}
	ctx.Cancel()

	mainHand, _ := ctx.Val().HeldItems()
	lobbyItemIndex, ok := mainHand.Value("lobby_item")
	if !ok {
		return
	}

	switch lobbyItemIndex {
	case 0: // play ffa
		sendFFAForm(ctx.Val())
	case 1: // play duels
	case 8: // profile
	}
}
