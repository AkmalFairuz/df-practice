package ffa

import (
	"github.com/df-mc/dragonfly/server/entity"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/mcdb"
	"log/slog"
)

func InitArenas(log *slog.Logger) {
	initClassicArena(log)
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
