package ffa

import (
	"github.com/akmalfairuz/df-practice/practice/helper"
	"github.com/akmalfairuz/df-practice/practice/kit"
	"log/slog"
)

var noDebuffArena *Arena

func NoDebuffArena() *Arena {
	return noDebuffArena
}

func initNoDebuffArena(log *slog.Logger) {
	noDebuffArena = New(initWorldConfig(log, "worlds/ffa_no_debuff").New())
	noDebuffArena.voidY = 70
	noDebuffArena.spawns = []helper.Location{
		{X: 46.5, Y: 76.1, Z: 44.5, Yaw: 135, Pitch: 0},
		{X: 48.5, Y: 76.1, Z: 44.5, Yaw: 135, Pitch: 0},
		{X: 46.5, Y: 76.1, Z: 46.5, Yaw: 135, Pitch: 0},
		{X: 48.5, Y: 76.1, Z: 46.5, Yaw: 135, Pitch: 0},
		{X: 50.5, Y: 76.1, Z: 44.5, Yaw: 135, Pitch: 0},
	}
	noDebuffArena.icon = "textures/items/potion_bottle_splash_heal"
	noDebuffArena.k = kit.NoDebuff{}
	if err := noDebuffArena.Init(); err != nil {
		panic(err)
	}
}
