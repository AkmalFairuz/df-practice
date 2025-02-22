package practice

import (
	"github.com/akmalfairuz/df-practice/practice/user"
	"github.com/df-mc/dragonfly/server/event"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
	"strings"
)

type packetHandler struct {
}

func (p packetHandler) HandleClientPacket(ctx *event.Context[*player.Player], pk packet.Packet) {
	switch pk := pk.(type) {
	case *packet.ClientCacheStatus:
		pk.Enabled = false
	case *packet.PlayerAuthInput:
		if pk.InputData.Load(packet.InputFlagMissedSwing) {
			addUserClick(ctx.Val())
		}
	case *packet.InventoryTransaction:
		trData, ok := pk.TransactionData.(*protocol.UseItemOnEntityTransactionData)
		if ok {
			if trData.ActionType == protocol.UseItemOnEntityActionAttack {
				addUserClick(ctx.Val())
			}
		}
	case *packet.LevelSoundEvent:
		if pk.SoundType == packet.SoundEventAttackNoDamage {
			addUserClick(ctx.Val())
		}
	case *packet.Text:
		// Execute command if it starts with "./" like PocketMine does
		if strings.HasPrefix(pk.Message, "./") {
			ctx.Val().ExecuteCommand(pk.Message[1:])
		}
	}
}

func addUserClick(p *player.Player) {
	u := user.Get(p)
	if u == nil {
		return
	}
	u.HandleClientClick()
}

func (p packetHandler) HandleServerPacket(ctx *event.Context[*player.Player], pk packet.Packet) {
	switch pk := pk.(type) {
	case *packet.LevelSoundEvent:
		if pk.SoundType == packet.SoundEventAttackStrong || pk.SoundType == packet.SoundEventAttackNoDamage || pk.SoundType == packet.SoundEventAttack {
			ctx.Cancel() // disable hit sound
		}
	}
}
