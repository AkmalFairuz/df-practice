package igame

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/item/inventory"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
	"time"
)

type IGame interface {
	Join(p *player.Player) error
	Quit(p *player.Player) error
	SetCloseHook(func())
	ID() string
	IsWaiting() bool
	IsPlaying() bool
	IsEnding() bool
	CurrentTick() uint64
	Participants() map[string]IParticipant
	PlayingParticipants() map[string]IParticipant
	End()
	World() *world.World
	Players(tx *world.Tx) []*player.Player
	SetSpectator(p *player.Player)
	HandleBlockBreak(ctx *player.Context, pos cube.Pos, drops *[]item.Stack, xp *int)
	HandleBlockPlace(ctx *player.Context, pos cube.Pos, b world.Block)
	HandleMove(ctx *player.Context, pos mgl64.Vec3, rot cube.Rotation)
	HandleFoodLoss(ctx *player.Context, from int, to *int)
	HandleHeal(ctx *player.Context, health *float64, src world.HealingSource)
	HandleHurt(ctx *player.Context, damage *float64, immune bool, immunity *time.Duration, src world.DamageSource)
	HandleItemUse(ctx *player.Context)
	HandleItemUseOnEntity(ctx *player.Context, e world.Entity)
	HandleAttackEntity(ctx *player.Context, e world.Entity, force *float64, height *float64, critical *bool)
	HandleDrop(ctx *inventory.Context, slot int, stack item.Stack)
	HandlePlace(ctx *inventory.Context, slot int, stack item.Stack)
	HandleTake(ctx *inventory.Context, slot int, stack item.Stack)
	Load() error
}
