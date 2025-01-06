package command

import (
	"github.com/akmalfairuz/df-practice/practice/user"
	"github.com/akmalfairuz/df-practice/translations"
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
	onlyAdmin

	Type    gameModeEnum               `cmd:"gamemode"`
	Targets cmd.Optional[onlineTarget] `cmd:"target"`
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
	target := g.Targets.LoadOr(onlineTarget{name: s.(*player.Player).Name(), u: user.Get(s.(*player.Player))})

	target.ExecutePlayer(func(p *player.Player, ok bool) {
		p.SetGameMode(getGameModeFromString(string(g.Type)))
		messaget(s, translations.GamemodeCommandSuccess, p.Name(), g.Type)
	})
}
