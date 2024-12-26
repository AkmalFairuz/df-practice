package main

import (
	"github.com/akmalfairuz/df-practice/practice"
	"github.com/akmalfairuz/df-practice/practice/command"
	"github.com/akmalfairuz/df-practice/practice/config"
	"github.com/akmalfairuz/df-practice/practice/game"
	"github.com/df-mc/dragonfly/server"
	"github.com/df-mc/dragonfly/server/cmd"
	"log/slog"
)

func main() {
	log := slog.Default()

	userConfig := server.DefaultConfig()
	userConfig.Server.Name = config.Get().Server.Name
	userConfig.Server.DisableJoinQuitMessages = true

	userConfig.Players.SaveData = false
	userConfig.Players.MaxCount = config.Get().Server.MaxPlayers
	userConfig.Network.Address = config.Get().Server.ListenAddress

	serverConfig, err := userConfig.Config(log)
	if err != nil {
		panic(err)
	}

	cmd.Register(cmd.New("whisper", "", []string{"w", "msg"}, command.Whisper{}))
	cmd.Register(cmd.New("reply", "", []string{"r"}, command.Reply{}))
	cmd.Register(cmd.New("lobby", "", []string{"hub"}, command.Lobby{}))
	cmd.Register(cmd.New("gamemode", "", []string{"gm"}, command.GameMode{}))

	if err := game.Init(); err != nil {
		panic(err)
	}

	pr := practice.New(log, serverConfig.New())
	pr.Run()
}
