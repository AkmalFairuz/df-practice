package ffa

import (
	"github.com/akmalfairuz/df-practice/practice/kit"
	"log/slog"
)

var sumoArena *Arena

func initSumoArena(log *slog.Logger) {
	sumoArena = New(initWorldConfig(log, "worlds/ffa_sumo").New())
	sumoArena.applyConfig(configs["sumo"])
	sumoArena.zeroDamageExceptVoid = true
	sumoArena.k = kit.Nop{}
	sumoArena.disableHPNameTag = true
	sumoArena.disableHunger = true
	if err := sumoArena.Init(); err != nil {
		panic(err)
	}
}

func SumoArena() *Arena {
	return sumoArena
}
