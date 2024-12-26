package nop

import (
	"github.com/akmalfairuz/df-practice/practice/game"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/item/inventory"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
	"time"
)

type Game struct {
}

func (g Game) MinimumParticipants() int {
	panic("must be implemented")
}

func (g Game) PlayingTime() int {
	panic("must be implemented")
}

func (g Game) WaitingTime() int {
	panic("must be implemented")
}

func (g Game) EndingTime() int {
	panic("must be implemented")
}

func (g Game) MaxParticipants() int {
	panic("must be implemented")
}

func (g Game) CheckEnd() {
	
}

func (g Game) GameName() string {
	return "Nop"
}

func (g Game) OnJoin(p *player.Player) error {
	return nil
}

func (g Game) OnJoined(par *game.Participant, p *player.Player) {
}

func (g Game) OnQuit(p *player.Player) {
}

func (g Game) OnInit() {
}

func (g Game) OnStart() {
}

func (g Game) OnEnd() {
}

func (g Game) OnStop() {
}

func (g Game) OnTick() {
}

func (g Game) HandleHurt(ctx *player.Context, damage *float64, immune bool, immunity *time.Duration, src world.DamageSource) {
}

func (g Game) HandleHeal(ctx *player.Context, health *float64, src world.HealingSource) {
}

func (g Game) HandleFoodLoss(ctx *player.Context, from int, to *int) {
}

func (g Game) HandleBlockBreak(ctx *player.Context, pos cube.Pos, drops *[]item.Stack, xp *int) {
}

func (g Game) HandleBlockPlace(ctx *player.Context, pos cube.Pos, b world.Block) {
}

func (g Game) HandleMove(ctx *player.Context, pos mgl64.Vec3, rot cube.Rotation) {
}

func (g Game) HandleAttackEntity(ctx *player.Context, e world.Entity, force *float64, height *float64, critical *bool) {
}

func (g Game) HandleItemUse(ctx *player.Context) {
}

func (g Game) HandleItemUseOnEntity(ctx *player.Context, e world.Entity) {
}

func (g Game) HandleDrop(ctx *inventory.Context, slot int, stack item.Stack) {
}

func (g Game) HandlePlace(ctx *inventory.Context, slot int, stack item.Stack) {
}

func (g Game) HandleTake(ctx *inventory.Context, slot int, stack item.Stack) {
}

var _ game.Impl = (*Game)(nil)
