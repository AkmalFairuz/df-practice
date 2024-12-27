package kit

import (
	"github.com/akmalfairuz/df-practice/practice/helper"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/item/enchantment"
	"github.com/df-mc/dragonfly/server/item/potion"
	"github.com/df-mc/dragonfly/server/player"
)

type Classic struct{}

func (Classic) Items(*player.Player) [36]item.Stack {
	items := [36]item.Stack{
		helper.SetItemAsUnbreakable(item.NewStack(item.Sword{Tier: item.ToolTierIron}, 1)),
		helper.SetItemAsUnbreakable(item.NewStack(item.Bow{}, 1)),
		item.NewStack(item.GoldenApple{}, 16),
		item.NewStack(item.Potion{Type: potion.Healing()}, 5),
		item.NewStack(item.Arrow{}, 16),
	}
	return items
}

func (Classic) Armour(*player.Player) [4]item.Stack {
	protection := item.NewEnchantment(enchantment.Protection, 3)
	return [4]item.Stack{
		helper.SetItemAsUnbreakable(item.NewStack(item.Helmet{Tier: item.ArmourTierChain{}}, 1).WithEnchantments(protection)),
		helper.SetItemAsUnbreakable(item.NewStack(item.Chestplate{Tier: item.ArmourTierChain{}}, 1).WithEnchantments(protection)),
		helper.SetItemAsUnbreakable(item.NewStack(item.Leggings{Tier: item.ArmourTierChain{}}, 1).WithEnchantments(protection)),
		helper.SetItemAsUnbreakable(item.NewStack(item.Boots{Tier: item.ArmourTierChain{}}, 1).WithEnchantments(protection)),
	}
}
