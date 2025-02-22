package lobby

import (
	"fmt"
	"github.com/akmalfairuz/df-practice/practice/helper"
	"github.com/akmalfairuz/df-practice/practice/lang"
	"github.com/akmalfairuz/df-practice/practice/user"
	"github.com/akmalfairuz/df-practice/translations"
	"github.com/df-mc/dragonfly/server/block/cube"
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
	w.SetSpawn(cube.Pos{0, 69, 0})
	instance = ret
	initFFAEntries()
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
		u.Translatef(translations.LobbyScoreboardPlayersOnlineCount, user.Count()),
		u.Translatef(translations.LobbyScoreboardPlayersPlayingCount, user.Count()-len(l.users(tx))),
		"",
		u.Translatef(translations.ScoreboardYourPing, u.Session().Latency().Milliseconds()),
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
	if u == nil {
		return
	}

	_ = p.Inventory().SetItem(0, item.NewStack(item.Sword{Tier: item.ToolTierDiamond}, 1).WithCustomName(lang.Translatef(u.Lang(), translations.LobbyItemPlayFfaName)).WithValue(lobbyItemIndexKey, 0))
	_ = p.Inventory().SetItem(1, item.NewStack(item.Sword{Tier: item.ToolTierIron}, 1).WithCustomName(lang.Translatef(u.Lang(), translations.LobbyItemPlayDuelsName)).WithValue(lobbyItemIndexKey, 1))
	//_ = p.Inventory().SetItem(8, item.NewStack(block.Skull{Type: block.PlayerHead()}, 1).WithCustomName(lang.Translatef(u.Lang(), translations.LobbyItemProfileName)).WithValue(lobbyItemIndexKey, 8))
}

func (l *Lobby) Spawn(p *player.Player) {
	l.doSpawn(p, false)
}

func (l *Lobby) SpawnSync(p *player.Player) {
	l.doSpawn(p, true)
}

func (l *Lobby) doSpawn(p *player.Player, sync bool) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("panic recovered when spawning player to lobby: %v\n", r)
		}
	}()

	if p.Tx().World() == l.w {
		p.Teleport(l.w.Spawn().Vec3())
		helper.ResetPlayer(p)
		l.sendLobbyItems(p)
		l.sendUserScoreboard(user.Get(p), p.Tx())
	} else {
		ent, _ := p.H().Entity(p.Tx())

		p.Tx().RemoveEntity(ent)

		exec := l.w.Exec(func(tx *world.Tx) {
			newP := tx.AddEntity(p.H()).(*player.Player)
			helper.ResetPlayer(newP)
			u := user.Get(newP)
			if u != nil {
				l.sendLobbyItems(newP)
				l.sendUserScoreboard(u, tx)
			}
			newP.SetGameMode(world.GameModeAdventure)
			newP.Teleport(l.w.Spawn().Vec3())
		})
		if sync {
			<-exec
		}
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
	case 1:
		sendDuelsForm(ctx.Val())
	case 8: // profile
	}
}
