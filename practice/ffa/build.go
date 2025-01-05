package ffa

import (
	"github.com/akmalfairuz/df-practice/practice/kit"
	"log/slog"
)

var buildArena *Arena

func BuildArena() *Arena {
	return buildArena
}

func initBuildArena(log *slog.Logger) {
	buildArena = New(initWorldConfig(log, "worlds/ffa_build").New())
	buildArena.applyConfig(configs["build"])
	buildArena.disableHunger = true
	buildArena.k = kit.Build{}
	buildArena.allowBuild = true
	if err := buildArena.Init(); err != nil {
		panic(err)
	}
}
