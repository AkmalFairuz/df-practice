package kit

import (
	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/item/enchantment"
	"github.com/df-mc/dragonfly/server/player"
)

type Build struct{}

func (b Build) Items(*player.Player) [36]item.Stack {
	durability := item.NewEnchantment(enchantment.Unbreaking, 10)
	return [36]item.Stack{
		item.NewStack(item.Sword{Tier: item.ToolTierStone}, 1).WithEnchantments(durability),
		item.NewStack(item.Pickaxe{Tier: item.ToolTierIron}, 1).WithEnchantments(durability),
		item.NewStack(block.CoralBlock{Type: block.BrainCoral(), Dead: true}, 64),
		item.NewStack(block.CoralBlock{Type: block.BrainCoral(), Dead: true}, 64),
		item.NewStack(item.GoldenApple{}, 5),
	}
}

func (b Build) Armour(*player.Player) [4]item.Stack {
	protection := item.NewEnchantment(enchantment.Protection, 1)
	durability := item.NewEnchantment(enchantment.Unbreaking, 10)
	return [4]item.Stack{
		{},
		item.NewStack(item.Chestplate{Tier: item.ArmourTierIron{}}, 1).WithEnchantments(protection, durability),
		item.NewStack(item.Leggings{Tier: item.ArmourTierIron{}}, 1).WithEnchantments(protection, durability),
		item.NewStack(item.Boots{Tier: item.ArmourTierIron{}}, 1).WithEnchantments(protection, durability),
	}
}
