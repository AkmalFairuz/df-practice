package helper

import (
	"github.com/df-mc/dragonfly/server/player"
	"github.com/sandertv/gophertunnel/minecraft/text"
	"time"
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
	p.Extinguish()
}

func UpdatePlayerNameTagWithHealth(p *player.Player, additionalHealth float64) {
	totalHealth := max(min(p.Health()+additionalHealth, p.MaxHealth()), 0)
	p.SetNameTag(p.Name() + "\n" + text.Colourf("<red>❤</red><white>%.0f</white>", totalHealth))
}

func UpdateXPBarCooldownDisplay(p *player.Player, lastCooldownStart time.Time, cooldown time.Duration) {
	if time.Since(lastCooldownStart) < cooldown {
		p.SetExperienceLevel(int(time.Until(lastCooldownStart.Add(cooldown)).Seconds()) + 1)
		p.SetExperienceProgress(1 - float64(time.Since(lastCooldownStart))/float64(cooldown))
	} else if p.ExperienceLevel() != 0 || p.ExperienceProgress() >= 0.01 {
		p.SetExperienceLevel(0)
		p.SetExperienceProgress(0)
	}
}
