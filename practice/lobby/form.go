package lobby

import (
	"github.com/akmalfairuz/df-practice/practice/ffa"
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
					_ = ffa.ClassicArena().Join(p, tx)
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
