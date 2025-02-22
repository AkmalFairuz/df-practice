package duelsmanager

import (
	"github.com/akmalfairuz/df-practice/practice/game"
	"github.com/akmalfairuz/df-practice/practice/game/duels"
	"github.com/akmalfairuz/df-practice/practice/game/gamemanager"
	"github.com/akmalfairuz/df-practice/practice/game/igame"
	"github.com/akmalfairuz/df-practice/practice/kit"
	"log/slog"
)

var Classic = gamemanager.New(func(mgr *gamemanager.Manager) igame.Impl {
	g := ClassicRaw()
	g.Game().SetPlayAgainHook(mgr.Join)
	return g
})

var ClassicRaw = func() igame.Impl {
	d := &duels.Duels{}
	d.SetKit(kit.Classic{})

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
}

var NoDebuff = gamemanager.New(func(mgr *gamemanager.Manager) igame.Impl {
	g := NoDebuffRaw()
	g.Game().SetPlayAgainHook(mgr.Join)
	return g
})

var NoDebuffRaw = func() igame.Impl {
	d := &duels.Duels{}
	d.SetKit(kit.NoDebuff{})

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
}
