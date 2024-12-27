package lobby

import (
	"github.com/akmalfairuz/df-practice/practice/ffa"
	"github.com/akmalfairuz/df-practice/practice/game/duelsmanager"
	"github.com/akmalfairuz/df-practice/practice/game/gamemanager"
	"github.com/akmalfairuz/df-practice/practice/helper"
	"github.com/akmalfairuz/df-practice/practice/user"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	form "github.com/twistedasylummc/inline-forms"
)

func sendFFAForm(p *player.Player) {
	u := user.Get(p)

	p.SendForm(&form.Menu{
		Title: u.Translatef("form.ffa.selector.title"),
		Buttons: []form.Button{
			{
				Text:  u.Translatef("form.ffa.selector.classic"),
				Image: ffa.ClassicArena().Icon(),
				Submit: func(tx *world.Tx) {
					helper.LogErrors(ffa.ClassicArena().Join(p, tx))
				},
			},
			//{
			//	Text: "NoDebuff",
			//	Submit: func(tx *world.Tx) {
			//		_ = ffa.NoDebuffArena().Join(p, tx)
			//	},
			//},
			//{
			//	Text: "Build",
			//	Submit: func(tx *world.Tx) {
			//	},
			//},
		},
	})
}

type duelsEntry struct {
	TranslationName string
	Manager         *gamemanager.Manager
}

var duelsEntries = []duelsEntry{
	{
		TranslationName: "form.duels.selector.classic",
		Manager:         duelsmanager.Classic,
	},
	{
		TranslationName: "form.duels.selector.nodebuff",
		Manager:         duelsmanager.NoDebuff,
	},
}

func sendDuelsForm(p *player.Player) {
	u := user.Get(p)

	btns := make([]form.Button, 0, len(duelsEntries))

	for _, entry := range duelsEntries {
		btns = append(btns, form.Button{
			Text: u.Translatef(entry.TranslationName) + "\n" + u.Translatef("form.playing.format", entry.Manager.PlayersCount()),
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
		Title:   u.Translatef("form.duels.selector.title"),
		Buttons: btns,
	})
}
