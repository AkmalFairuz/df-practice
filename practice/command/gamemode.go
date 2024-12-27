package command

import (
	"github.com/akmalfairuz/df-practice/practice/user"
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
)

type gameModeEnum string

func (gameModeEnum) Type() string {
	return "GameMode"
}

func (gameModeEnum) Options(_ cmd.Source) []string {
	return []string{"survival", "creative", "adventure", "spectator", "s", "c", "a", "sp"}
}

type GameMode struct {
	Type    gameModeEnum               `cmd:"gamemode"`
	Targets cmd.Optional[[]cmd.Target] `cmd:"target"`
}

func getGameModeFromString(str string) world.GameMode {
	switch str {
	case "creative", "c":
		return world.GameModeCreative
	case "adventure", "a":
		return world.GameModeAdventure
	case "spectator", "sp":
		return world.GameModeSpectator
	default:
		return world.GameModeSurvival
	}
}

func (g GameMode) Run(s cmd.Source, o *cmd.Output, tx *world.Tx) {
	targets := g.Targets.LoadOr([]cmd.Target{s})

	for _, t := range targets {
		if p, ok := t.(*player.Player); ok {
			p.SetGameMode(getGameModeFromString(string(g.Type)))
			o.Printf(translatef(s, "gamemode.command.success", p.Name(), g.Type))
		}
	}
}

func (g GameMode) Allow(s cmd.Source) bool {
	if p, ok := s.(*player.Player); ok {
		u := user.Get(p)
		if u == nil {
			return false
		}
		return u.RankName() == "admin"
	}
	return false
}
