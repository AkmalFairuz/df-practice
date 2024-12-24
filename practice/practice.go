package practice

import (
	"fmt"
	"github.com/akmalfairuz/df-practice/practice/ffa"
	"github.com/akmalfairuz/df-practice/practice/lang"
	"github.com/akmalfairuz/df-practice/practice/lobby"
	"github.com/akmalfairuz/df-practice/practice/user"
	"github.com/bedrock-gophers/intercept/intercept"
	"github.com/df-mc/dragonfly/server"
	"github.com/df-mc/dragonfly/server/world"
	"log/slog"
)

type Practice struct {
	log *slog.Logger
	srv *server.Server

	l *lobby.Lobby
}

func New(log *slog.Logger, srv *server.Server) *Practice {
	l := lobby.New(srv.World())
	ffa.InitArenas(log)

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

		u := user.New(p)
		if err := u.Load(); err != nil {
			pr.log.Error("failed to load user data", "error", err)
			u.Disconnect(lang.Translate(u.Lang(), "user.load.error"))
			continue
		}
		user.Store(u)

		intercept.Intercept(p)
		p.Handle(newPlayerHandler(pr, u))
		p.Inventory().Handle(newPlayerInvHandler(p.Inventory()))
		p.Inventory().Handle(newPlayerInvHandler(p.Armour().Inventory()))

		user.BroadcastMessaget("player.join.message", p.Name())

		if err := pr.l.Init(); err != nil {
			panic(fmt.Errorf("failed to init lobby: %w", err))
		}
		pr.l.Spawn(p)

		go startPlayerTick(u)
	}
}
