package game

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/item/inventory"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
	"time"
)

type Impl interface {
	// PlayingTime returns the playing time of the game in seconds.
	PlayingTime() int
	// WaitingTime returns the waiting time of the game in seconds.
	WaitingTime() int
	// EndingTime returns the ending time of the game in seconds.
	EndingTime() int
	// MaxParticipants returns the maximum number of participants in the game.
	MaxParticipants() int
	// MinimumParticipants returns the minimum number of participants in the game.
	MinimumParticipants() int
	// GameName returns the name of the game.
	GameName() string
	// OnJoin is called when a player joins the game. Return an error to prevent the player from joining.
	OnJoin(p *player.Player) error
	// OnJoined is called when a player has successfully joined the game.
	OnJoined(par *Participant, p *player.Player)
	// OnQuit is called when a player quits the game.
	OnQuit(p *player.Player)
	// OnInit is called when the game is initialized.
	OnInit()
	// OnStart is called when the game starts.
	OnStart()
	// OnEnd is called when the game ends.
	OnEnd()
	// OnStop is called when the game stops.
	OnStop()
	// OnTick is called every 50 milliseconds.
	OnTick()
	// CheckEnd is called when the game should check if it should end or not.
	CheckEnd()

	HandleHurt(ctx *player.Context, damage *float64, immune bool, immunity *time.Duration, src world.DamageSource)
	HandleHeal(ctx *player.Context, health *float64, src world.HealingSource)
	HandleFoodLoss(ctx *player.Context, from int, to *int)
	HandleBlockBreak(ctx *player.Context, pos cube.Pos, drops *[]item.Stack, xp *int)
	HandleBlockPlace(ctx *player.Context, pos cube.Pos, b world.Block)
	HandleMove(ctx *player.Context, pos mgl64.Vec3, rot cube.Rotation)
	HandleAttackEntity(ctx *player.Context, e world.Entity, force *float64, height *float64, critical *bool)
	HandleItemUse(ctx *player.Context)
	HandleItemUseOnEntity(ctx *player.Context, e world.Entity)
	HandleDrop(ctx *inventory.Context, slot int, stack item.Stack)
	HandlePlace(ctx *inventory.Context, slot int, stack item.Stack)
	HandleTake(ctx *inventory.Context, slot int, stack item.Stack)
}

type ParticipantImpl interface {
}
