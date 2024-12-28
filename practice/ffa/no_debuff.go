package ffa

import (
	"github.com/akmalfairuz/df-practice/practice/kit"
	"log/slog"
)

var noDebuffArena *Arena

func NoDebuffArena() *Arena {
	return noDebuffArena
}

func initNoDebuffArena(log *slog.Logger) {
	noDebuffArena = New(initWorldConfig(log, "worlds/ffa_no_debuff").New())
	noDebuffArena.applyConfig(configs["no_debuff"])
	noDebuffArena.k = kit.NoDebuff{}
	if err := noDebuffArena.Init(); err != nil {
		panic(err)
	}
}
