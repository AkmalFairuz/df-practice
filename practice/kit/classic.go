package kit

import (
	"github.com/akmalfairuz/df-practice/practice/helper"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/item/enchantment"
	"github.com/df-mc/dragonfly/server/item/potion"
	"github.com/df-mc/dragonfly/server/player"
)

func Classic(p *player.Player) error {
	_, _ = p.Inventory().AddItem(helper.SetItemAsUnbreakable(item.NewStack(item.Sword{Tier: item.ToolTierIron}, 1)))
	_, _ = p.Inventory().AddItem(helper.SetItemAsUnbreakable(item.NewStack(item.Bow{}, 1)))
	_, _ = p.Inventory().AddItem(item.NewStack(item.GoldenApple{}, 16))
	_, _ = p.Inventory().AddItem(item.NewStack(item.Potion{Type: potion.Healing()}, 5))
	_, _ = p.Inventory().AddItem(item.NewStack(item.Arrow{}, 16))

	p.Armour().SetHelmet(helper.SetItemAsUnbreakable(item.NewStack(item.Helmet{Tier: item.ArmourTierChain{}}, 1).WithEnchantments(item.NewEnchantment(enchantment.Protection, 3))))
	p.Armour().SetChestplate(helper.SetItemAsUnbreakable(item.NewStack(item.Chestplate{Tier: item.ArmourTierChain{}}, 1).WithEnchantments(item.NewEnchantment(enchantment.Protection, 3))))
	p.Armour().SetLeggings(helper.SetItemAsUnbreakable(item.NewStack(item.Leggings{Tier: item.ArmourTierChain{}}, 1).WithEnchantments(item.NewEnchantment(enchantment.Protection, 3))))
	p.Armour().SetBoots(helper.SetItemAsUnbreakable(item.NewStack(item.Boots{Tier: item.ArmourTierChain{}}, 1).WithEnchantments(item.NewEnchantment(enchantment.Protection, 3))))
	return nil
}
