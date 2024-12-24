package practice

import (
	"github.com/df-mc/dragonfly/server/event"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

type packetHandler struct {
}

func (p packetHandler) HandleClientPacket(ctx *event.Context[*player.Player], pk packet.Packet) {

}

func (p packetHandler) HandleServerPacket(ctx *event.Context[*player.Player], pk packet.Packet) {
	switch pk := pk.(type) {
	case *packet.LevelSoundEvent:
		if pk.SoundType == packet.SoundEventAttackStrong || pk.SoundType == packet.SoundEventAttackNoDamage {
			ctx.Cancel() // disable hit sound
		}
	}
}
