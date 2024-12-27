package ffa

import (
	"github.com/akmalfairuz/df-practice/practice/helper"
	"github.com/akmalfairuz/df-practice/practice/kit"
	"github.com/df-mc/dragonfly/server/entity"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/mcdb"
	"log/slog"
)

var classicArena *Arena
var noDebuffArena *Arena

func InitArenas(log *slog.Logger) {
	classicArena = New(initWorldConfig(log, "worlds/classic").New())
	// TODO: do not hardcode
	classicArena.spawns = []helper.Location{
		{X: -14.5, Y: 100.2, Z: -25.5, Yaw: -32, Pitch: 0},
		{X: -9.5, Y: 100.2, Z: -28.5, Yaw: -25, Pitch: 0},
		{X: 2.5, Y: 100.2, Z: -30.5, Yaw: 0, Pitch: 0},
		{X: 14.5, Y: 100.2, Z: -28.5, Yaw: 25, Pitch: 0},
		{X: 9.5, Y: 100.2, Z: 25.5, Yaw: 32, Pitch: 0},
	}
	classicArena.voidY = 80
	classicArena.icon = "textures/items/iron_sword.png"
	classicArena.k = kit.Classic{}
	if err := classicArena.Init(); err != nil {
		panic(err)
	}
}

func initWorldConfig(log *slog.Logger, path string) world.Config {
	prov, err := mcdb.Config{Log: log}.Open(path)
	if err != nil {
		panic(err)
	}
	return world.Config{
		Log:       log,
		Entities:  entity.DefaultRegistry,
		Dim:       world.Overworld,
		Provider:  prov,
		ReadOnly:  true,
		Generator: world.NopGenerator{},
	}
}

func ClassicArena() *Arena {
	return classicArena
}

func NoDebuffArena() *Arena {
	return noDebuffArena
}
