package kit

import (
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
)

type Nop struct{}

func (Nop) Items(*player.Player) [36]item.Stack {
	return [36]item.Stack{}
}

func (Nop) Armour(*player.Player) [4]item.Stack {
	return [4]item.Stack{}
}
