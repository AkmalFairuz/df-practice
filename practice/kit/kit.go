package kit

import (
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
)

type Kit interface {
	Items(p *player.Player) [36]item.Stack
	Armour(p *player.Player) [4]item.Stack
}

func Apply(k Kit, p *player.Player) {
	for i, it := range k.Items(p) {
		_ = p.Inventory().SetItem(i, it)
	}
	armour := k.Armour(p)
	p.Armour().Set(armour[0], armour[1], armour[2], armour[3])
}
