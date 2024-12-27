package gamedefaults

import (
	"github.com/akmalfairuz/df-practice/practice/game/igame"
	"github.com/df-mc/dragonfly/server/entity"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
)

func HandleKillMessage(g igame.IGame, p *player.Player, src world.DamageSource) {
	handled := false
	switch src := src.(type) {
	case entity.AttackDamageSource:
		att, ok := src.Attacker.(*player.Player)
		if !ok {
			break
		}

		handled = true
		g.Messaget("killed.message.format", p.Name(), att.Name())
	case entity.ProjectileDamageSource:
		att, ok := src.Owner.(*player.Player)
		if !ok {
			break
		}

		handled = true
		g.Messaget("killed.shot.message.format", p.Name(), att.Name())
	case entity.VoidDamageSource:
		// TODO: Handle void damage source
	}

	if !handled {
		g.Messaget("death.message.format", p.Name())
	}
}
