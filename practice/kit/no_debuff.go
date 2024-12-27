package kit

import (
	"github.com/akmalfairuz/df-practice/practice/kit/customitem"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/item/enchantment"
	"github.com/df-mc/dragonfly/server/item/potion"
	"github.com/df-mc/dragonfly/server/player"
)

type NoDebuff struct{}

func (NoDebuff) Items(*player.Player) [36]item.Stack {
	items := [36]item.Stack{}
	items[0] = item.NewStack(item.Sword{Tier: item.ToolTierDiamond}, 1).WithEnchantments(item.NewEnchantment(enchantment.Sharpness, 2), item.NewEnchantment(enchantment.FireAspect, 1), item.NewEnchantment(enchantment.Unbreaking, 10))
	items[1] = item.NewStack(customitem.NoDamageEnderPearl{}, 16)
	items[2] = item.NewStack(item.Beef{Cooked: true}, 64)
	for i := 3; i < 36; i++ {
		items[i] = item.NewStack(item.SplashPotion{Type: potion.StrongHealing()}, 1)
	}

	items[3] = item.NewStack(item.Potion{Type: potion.Swiftness()}, 1)
	items[8] = item.NewStack(item.Potion{Type: potion.FireResistance()}, 1)
	return items
}

func (NoDebuff) Armour(*player.Player) [4]item.Stack {
	durability := item.NewEnchantment(enchantment.Unbreaking, 10)
	protection := item.NewEnchantment(enchantment.Protection, 3)
	return [4]item.Stack{
		item.NewStack(item.Helmet{Tier: item.ArmourTierDiamond{}}, 1).WithEnchantments(durability, protection),
		item.NewStack(item.Chestplate{Tier: item.ArmourTierDiamond{}}, 1).WithEnchantments(durability, protection),
		item.NewStack(item.Leggings{Tier: item.ArmourTierDiamond{}}, 1).WithEnchantments(durability, protection),
		item.NewStack(item.Boots{Tier: item.ArmourTierDiamond{}}, 1).WithEnchantments(durability, protection),
	}
}
