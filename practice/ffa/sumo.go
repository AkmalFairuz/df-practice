package ffa

import (
	"github.com/akmalfairuz/df-practice/practice/helper"
	"github.com/akmalfairuz/df-practice/practice/kit"
	"log/slog"
)

var sumoArena *Arena

func initSumoArena(log *slog.Logger) {
	sumoArena = New(initWorldConfig(log, "worlds/ffa_sumo").New())
	sumoArena.voidY = 37
	sumoArena.spawns = []helper.Location{
		{X: 27.5, Y: 42.1, Z: -11.5, Yaw: 45},
		{X: 27.5, Y: 42.1, Z: 27.5, Yaw: 135},
		{X: -11.5, Y: 42.1, Z: 27.5, Yaw: -135},
		{X: -11.5, Y: 42.1, Z: -11.5, Yaw: -45},
	}
	sumoArena.zeroDamageExceptVoid = true
	sumoArena.icon = "textures/items/lead"
	sumoArena.k = kit.Nop{}
	sumoArena.attackCooldownTick = 8
	sumoArena.disableHPNameTag = true
	sumoArena.disableHunger = true
	if err := sumoArena.Init(); err != nil {
		panic(err)
	}
}

func SumoArena() *Arena {
	return sumoArena
}
