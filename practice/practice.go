package practice

import (
	"fmt"
	"github.com/akmalfairuz/df-practice/internal/meta"
	"github.com/akmalfairuz/df-practice/practice/ffa"
	"github.com/akmalfairuz/df-practice/practice/lang"
	"github.com/akmalfairuz/df-practice/practice/lobby"
	"github.com/akmalfairuz/df-practice/practice/user"
	"github.com/akmalfairuz/df-practice/translations"
	"github.com/bedrock-gophers/intercept/intercept"
	"github.com/df-mc/dragonfly/server"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"log/slog"
	"strconv"
	"time"
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

	go func() {
		broadcastNo := 0
		ticker := time.NewTicker(time.Second * 60)
		for {
			select {
			case <-ticker.C:
				broadcastNo++
				if broadcastNo >= 3 {
					broadcastNo = 1
				}
				user.BroadcastMessaget("broadcast." + strconv.Itoa(broadcastNo) + ".message")
			}
		}
	}()

	pr.srv.Listen()
	for p := range pr.srv.Accept() {
		pr.log.Info("player connected", "player", p.Name())

		p.SetGameMode(world.GameModeAdventure)

		go func() {
			u := user.New(p)

			//field := reflect.ValueOf(u.Conn()).Elem().FieldByName("cacheEnabled")
			//reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).Elem().SetBool(false)

			if err := u.Load(); err != nil {
				pr.log.Error("failed to load user data", "error", err)
				u.Disconnect(lang.Translate(u.Lang(), translations.UserLoadError))
				return
			}

			p.H().ExecWorld(func(tx *world.Tx, e world.Entity) {
				user.Store(u)
				u.SetWorld(tx.World())
				newP := e.(*player.Player)

				intercept.Intercept(newP)
				newP.Handle(newPlayerHandler(pr, u))
				newP.Inventory().Handle(newPlayerInvHandler(newP.Inventory()))
				newP.Inventory().Handle(newPlayerInvHandler(newP.Armour().Inventory()))

				user.BroadcastMessaget(translations.PlayerJoinMessage, newP.Name())

				if err := pr.l.Init(); err != nil {
					panic(fmt.Errorf("failed to init lobby: %w", err))
				}
				pr.l.Spawn(newP)

				u.Messaget(translations.WelcomeMessage, newP.Name())

				go startPlayerTick(u)
			})
		}()
	}
}
