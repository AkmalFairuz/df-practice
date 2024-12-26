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

type Nop struct {
}

func (g Nop) MinimumParticipants() int {
	panic("must be implemented")
}

func (g Nop) PlayingTime() int {
	panic("must be implemented")
}

func (g Nop) WaitingTime() int {
	panic("must be implemented")
}

func (g Nop) EndingTime() int {
	panic("must be implemented")
}

func (g Nop) MaxParticipants() int {
	panic("must be implemented")
}

func (g Nop) CheckEnd() {

}

func (g Nop) GameName() string {
	return "Nop"
}

func (g Nop) Create(_ IGame) {
	panic("must be implemented")
}

func (g Nop) OnJoin(p *player.Player) error {
	return nil
}

func (g Nop) OnJoined(par IParticipant, p *player.Player) {
}

func (g Nop) OnQuit(p *player.Player) {
}

func (g Nop) OnInit() {
}

func (g Nop) OnStart() {
}

func (g Nop) OnEnd() {
}

func (g Nop) OnStop() {
}

func (g Nop) OnTick() {
}

func (g Nop) Game() IGame {
	panic("must be implemented")
}

func (g Nop) HandleHurt(ctx *player.Context, damage *float64, immune bool, immunity *time.Duration, src world.DamageSource) {
}

func (g Nop) HandleHeal(ctx *player.Context, health *float64, src world.HealingSource) {
}

func (g Nop) HandleFoodLoss(ctx *player.Context, from int, to *int) {
}

func (g Nop) HandleBlockBreak(ctx *player.Context, pos cube.Pos, drops *[]item.Stack, xp *int) {
}

func (g Nop) HandleBlockPlace(ctx *player.Context, pos cube.Pos, b world.Block) {
}

func (g Nop) HandleMove(ctx *player.Context, pos mgl64.Vec3, rot cube.Rotation) {
}

func (g Nop) HandleAttackEntity(ctx *player.Context, e world.Entity, force *float64, height *float64, critical *bool) {
}

func (g Nop) HandleItemUse(ctx *player.Context) {
}

func (g Nop) HandleItemUseOnEntity(ctx *player.Context, e world.Entity) {
}

func (g Nop) HandleDrop(ctx *inventory.Context, slot int, stack item.Stack) {
}

func (g Nop) HandlePlace(ctx *inventory.Context, slot int, stack item.Stack) {
}

func (g Nop) HandleTake(ctx *inventory.Context, slot int, stack item.Stack) {
}

var _ Impl = (*Nop)(nil)
