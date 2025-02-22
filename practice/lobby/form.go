package lobby

import (
	"github.com/akmalfairuz/df-practice/practice/ffa"
	"github.com/akmalfairuz/df-practice/practice/game/duelsmanager"
	"github.com/akmalfairuz/df-practice/practice/game/gamemanager"
	"github.com/akmalfairuz/df-practice/practice/helper"
	"github.com/akmalfairuz/df-practice/practice/user"
	"github.com/akmalfairuz/df-practice/translations"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	form "github.com/twistedasylummc/inline-forms"
)

type ffaEntry struct {
	TranslationName string
	Arena           *ffa.Arena
}

var ffaEntries []ffaEntry

func initFFAEntries() {
	ffaEntries = []ffaEntry{
		{
			TranslationName: translations.FormFfaSelectorClassic,
			Arena:           ffa.ClassicArena(),
		},
		{
			TranslationName: translations.FormFfaSelectorNodebuff,
			Arena:           ffa.NoDebuffArena(),
		},
		{
			TranslationName: translations.FormFfaSelectorSumo,
			Arena:           ffa.SumoArena(),
		},
		{
			TranslationName: translations.FormFfaSelectorBuild,
			Arena:           ffa.BuildArena(),
		},
	}
}

func sendFFAForm(p *player.Player) {
	u := user.Get(p)

	btns := make([]form.Button, 0, len(ffaEntries))
	for _, entry := range ffaEntries {
		btns = append(btns, form.Button{
			Text: u.Translatef(entry.TranslationName) + "\n" + u.Translatef(translations.FormPlayingFormat, len(entry.Arena.Participants())),
			Submit: func(tx *world.Tx) {
				ent, _ := p.H().Entity(tx)
				ent2 := ent.(*player.Player)
				helper.LogErrors(entry.Arena.Join(ent2, ent2.Tx()))
			},
			Image: entry.Arena.Icon(),
		})
	}

	p.SendForm(&form.Menu{
		Title:   u.Translatef(translations.FormFfaSelectorTitle),
		Buttons: btns,
	})
}

type duelsEntry struct {
	TranslationName string
	Manager         *gamemanager.Manager
	Icon            string
}

var duelsEntries = []duelsEntry{
	{
		TranslationName: translations.FormDuelsSelectorClassic,
		Manager:         duelsmanager.Classic,
		Icon:            "textures/items/iron_sword",
	},
	{
		TranslationName: translations.FormDuelsSelectorNodebuff,
		Manager:         duelsmanager.NoDebuff,
		Icon:            "textures/items/potion_bottle_splash_heal",
	},
}

func sendDuelsForm(p *player.Player) {
	u := user.Get(p)

	btns := make([]form.Button, 0, len(duelsEntries))

	for _, entry := range duelsEntries {
		btns = append(btns, form.Button{
			Text: u.Translatef(entry.TranslationName) + "\n" + u.Translatef(translations.FormPlayingFormat, entry.Manager.PlayersCount()),
			Submit: func(tx *world.Tx) {
				ent, ok := p.H().Entity(tx)
				if !ok {
					return
				}
				helper.LogErrors(entry.Manager.Join(ent.(*player.Player)))
			},
		})
	}

	p.SendForm(&form.Menu{
		Title:   u.Translatef(translations.FormDuelsSelectorTitle),
		Buttons: btns,
	})
}
