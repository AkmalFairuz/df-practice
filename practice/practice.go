package practice

import (
	"fmt"
	"github.com/akmalfairuz/df-practice/internal/meta"
	"github.com/akmalfairuz/df-practice/practice/ffa"
	"github.com/akmalfairuz/df-practice/practice/lang"
	"github.com/akmalfairuz/df-practice/practice/lobby"
	"github.com/akmalfairuz/df-practice/practice/user"
	"github.com/bedrock-gophers/intercept/intercept"
	"github.com/df-mc/dragonfly/server"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"log/slog"
)

type Practice struct {
	log *slog.Logger
	srv *server.Server

	l *lobby.Lobby
}

func New(log *slog.Logger, srv *server.Server) *Practice {
	ffa.InitArenas(log)
	l := lobby.New(srv.World())

	meta.Set("lobby", lobby.Instance())

	return &Practice{
		log: log,
		srv: srv,
		l:   l,
	}
}

func (pr *Practice) Run() {
	intercept.Hook(&packetHandler{})

	pr.srv.Listen()
	for p := range pr.srv.Accept() {
		pr.log.Info("player connected", "player", p.Name())

		p.SetGameMode(world.GameModeAdventure)

		go func() {
			u := user.New(p)
			if err := u.Load(); err != nil {
				pr.log.Error("failed to load user data", "error", err)
				u.Disconnect(lang.Translate(u.Lang(), "user.load.error"))
				return
			}
			user.Store(u)

			p.H().ExecWorld(func(tx *world.Tx, e world.Entity) {
				newP := e.(*player.Player)

				intercept.Intercept(newP)
				newP.Handle(newPlayerHandler(pr, u))
				newP.Inventory().Handle(newPlayerInvHandler(newP.Inventory()))
				newP.Inventory().Handle(newPlayerInvHandler(newP.Armour().Inventory()))

				user.BroadcastMessaget("player.join.message", newP.Name())

				if err := pr.l.Init(); err != nil {
					panic(fmt.Errorf("failed to init lobby: %w", err))
				}
				pr.l.Spawn(newP)

				go startPlayerTick(u)
			})
		}()
	}
}
