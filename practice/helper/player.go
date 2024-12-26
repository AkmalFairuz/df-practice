package helper

import (
	"github.com/df-mc/dragonfly/server/player"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

type SetHealthSource struct {
}

func (SetHealthSource) HealingSource() {}

func ClearAllPlayerInv(p *player.Player) {
	p.Inventory().Clear()
	p.Armour().Clear()
	p.EnderChestInventory().Clear()
}

func ResetPlayerExp(p *player.Player) {
	p.SetExperienceLevel(0)
	p.SetExperienceProgress(0)
}

func ResetPlayerAttributes(p *player.Player) {
	for _, e := range p.Effects() {
		p.RemoveEffect(e.Type())
	}
	p.SetMaxHealth(20)
	p.SetFood(20)
	p.Heal(20, SetHealthSource{})
}

func ResetPlayer(p *player.Player) {
	//p.RemoveScoreboard()
	ClearAllPlayerInv(p)
	ResetPlayerExp(p)
	ResetPlayerAttributes(p)
	p.SetNameTag(p.Name())
}

func UpdatePlayerNameTagWithHealth(p *player.Player, additionalHealth float64) {
	totalHealth := max(min(p.Health()+additionalHealth, p.MaxHealth()), 0)
	p.SetNameTag(p.Name() + "\n" + text.Colourf("<red>❤</red><white>%.0f</white>", totalHealth))
}
