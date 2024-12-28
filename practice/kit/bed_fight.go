package kit

import (
	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/item/enchantment"
	"github.com/df-mc/dragonfly/server/player"
)

type BedFightTeamColor int

func (b BedFightTeamColor) Colour() item.Colour {
	switch b {
	case BedFightTeamRed:
		return item.ColourRed()
	case BedFightTeamBlue:
		return item.ColourBlue()
	}
	panic("invalid team colour")
}

const (
	BedFightTeamRed BedFightTeamColor = iota
	BedFightTeamBlue
)

type BedFight struct {
	TeamColor BedFightTeamColor
}

func (b BedFight) Items(*player.Player) [36]item.Stack {
	durability := item.NewEnchantment(enchantment.Unbreaking, 10)
	efficiency := item.NewEnchantment(enchantment.Efficiency, 1)

	items := [36]item.Stack{
		item.NewStack(item.Sword{Tier: item.ToolTierStone}, 1).WithEnchantments(durability),
		item.NewStack(block.Wool{Colour: b.TeamColor.Colour()}, 64),
		item.NewStack(item.Shears{}, 1).WithEnchantments(durability),
		item.NewStack(item.Pickaxe{Tier: item.ToolTierStone}, 1).WithEnchantments(durability, efficiency),
		item.NewStack(item.Axe{Tier: item.ToolTierStone}, 1).WithEnchantments(durability, efficiency),
	}
	return items
}

func (b BedFight) Armour(*player.Player) [4]item.Stack {
	durability := item.NewEnchantment(enchantment.Unbreaking, 10)
	protection := item.NewEnchantment(enchantment.Protection, 3)
	colour := b.TeamColor.Colour().RGBA()
	return [4]item.Stack{
		item.NewStack(item.Helmet{Tier: item.ArmourTierLeather{Colour: colour}}, 1).WithEnchantments(protection, durability),
		item.NewStack(item.Chestplate{Tier: item.ArmourTierLeather{Colour: colour}}, 1).WithEnchantments(protection, durability),
		item.NewStack(item.Leggings{Tier: item.ArmourTierLeather{Colour: colour}}, 1).WithEnchantments(protection, durability),
		item.NewStack(item.Boots{Tier: item.ArmourTierLeather{Colour: colour}}, 1).WithEnchantments(protection, durability),
	}
}
