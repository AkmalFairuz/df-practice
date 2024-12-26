package duelsmanager

import (
	"github.com/akmalfairuz/df-practice/practice/game"
	"github.com/akmalfairuz/df-practice/practice/game/duels"
	"github.com/akmalfairuz/df-practice/practice/game/gamemanager"
	"github.com/akmalfairuz/df-practice/practice/game/igame"
	"log/slog"
)

var Default = gamemanager.New(func() igame.Impl {
	d := &duels.Duels{}

	gConf := game.Config{
		Log:     slog.Default(),
		MapName: "classic",
		Impl:    d,
	}
	g, err := gConf.New()
	if err != nil {
		panic(err)
	}
	d.Create(g)
	return d
})
