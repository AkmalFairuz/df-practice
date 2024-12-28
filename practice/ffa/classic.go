package ffa

import (
	"github.com/akmalfairuz/df-practice/practice/kit"
	"log/slog"
)

var classicArena *Arena

func ClassicArena() *Arena {
	return classicArena
}

func initClassicArena(log *slog.Logger) {
	classicArena = New(initWorldConfig(log, "worlds/classic").New())
	classicArena.applyConfig(configs["classic"])
	classicArena.icon = "textures/items/iron_sword"
	classicArena.k = kit.Classic{}
	if err := classicArena.Init(); err != nil {
		panic(err)
	}
}
