package game

import (
	"errors"
	"github.com/akmalfairuz/df-practice/practice/helper"
	"github.com/df-mc/dragonfly/server/entity"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/mcdb"
	"log/slog"
	"path"
)

const (
	gameWorldsPath = "game_worlds"
)

type Config struct {
	Log *slog.Logger

	MapDir string

	Impl            Impl
	ParticipantImpl ParticipantImpl
}

func (c Config) New() (*Game, error) {
	if c.MapDir == "" {
		return nil, errors.New("game: map directory is required")
	}

	id := generateID()

	gameWorldPath := path.Join(gameWorldsPath, id)

	if err := helper.CopyDir(c.MapDir, gameWorldPath); err != nil {
		return nil, err
	}

	prov, err := mcdb.Open(gameWorldPath)
	if err != nil {
		return nil, err
	}

	var wConf world.Config
	wConf.Log = c.Log
	wConf.Dim = world.Overworld
	wConf.Provider = prov
	wConf.Entities = entity.DefaultRegistry
	w := wConf.New()
	w.SetDifficulty(world.DifficultyEasy)
	w.SetTime(3000)
	w.StopTime()
	w.StopThundering()
	w.StopRaining()
	w.StopWeatherCycle()

	return &Game{
		id:    id,
		log:   c.Log,
		w:     w,
		wDir:  gameWorldPath,
		impl:  c.Impl,
		pImpl: c.ParticipantImpl,
	}, nil
}
