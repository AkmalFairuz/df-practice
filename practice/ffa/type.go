package ffa

import (
	"github.com/akmalfairuz/df-practice/practice/helper"
	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/entity"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/item/enchantment"
	"github.com/df-mc/dragonfly/server/item/potion"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/biome"
	"github.com/df-mc/dragonfly/server/world/generator"
	"github.com/df-mc/dragonfly/server/world/mcdb"
	"log/slog"
)

var classicArena *Arena
var noDebuffArena *Arena

func InitArenas(log *slog.Logger) {
	classicArena = New(initWorldConfig(log, "worlds/classic").New())
	// TODO: do not hardcode
	classicArena.spawns = []helper.Location{
		{X: -14.5, Y: 100.2, Z: -25.5, Yaw: -32, Pitch: 0},
		{X: -9.5, Y: 100.2, Z: -28.5, Yaw: -25, Pitch: 0},
		{X: 2.5, Y: 100.2, Z: -30.5, Yaw: 0, Pitch: 0},
		{X: 14.5, Y: 100.2, Z: -28.5, Yaw: 25, Pitch: 0},
		{X: 9.5, Y: 100.2, Z: 25.5, Yaw: 32, Pitch: 0},
	}
	classicArena.voidY = 80
	classicArena.icon = "textures/items/iron_sword.png"
	classicArena.onSendKit = func(p *player.Player) error {
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
	if err := classicArena.Init(); err != nil {
		panic(err)
	}
}

func initWorldConfig(log *slog.Logger, path string) world.Config {
	prov, err := mcdb.Config{Log: log}.Open(path)
	if err != nil {
		panic(err)
	}
	return world.Config{
		Log:       log,
		Entities:  entity.DefaultRegistry,
		Dim:       world.Overworld,
		Provider:  prov,
		Generator: generator.NewFlat(biome.Plains{}, []world.Block{block.Grass{}, block.Dirt{}, block.Dirt{}, block.Bedrock{}}),
	}
}

func ClassicArena() *Arena {
	return classicArena
}

func NoDebuffArena() *Arena {
	return noDebuffArena
}
