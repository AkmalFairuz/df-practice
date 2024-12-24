package lobby

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
)

type lobbyWorldHandler struct {
	l *Lobby
}

func newLobbyWorldHandler(l *Lobby) *lobbyWorldHandler {
	return &lobbyWorldHandler{l: l}
}

func (l *lobbyWorldHandler) HandleLiquidFlow(ctx *world.Context, from, into cube.Pos, liquid world.Liquid, replaced world.Block) {
	ctx.Cancel()
}

func (l *lobbyWorldHandler) HandleLiquidDecay(ctx *world.Context, pos cube.Pos, before, after world.Liquid) {
	ctx.Cancel()
}

func (l *lobbyWorldHandler) HandleLiquidHarden(ctx *world.Context, hardenedPos cube.Pos, liquidHardened, otherLiquid, newBlock world.Block) {
	ctx.Cancel()
}

func (l *lobbyWorldHandler) HandleSound(ctx *world.Context, s world.Sound, pos mgl64.Vec3) {

}

func (l *lobbyWorldHandler) HandleFireSpread(ctx *world.Context, from, to cube.Pos) {
	ctx.Cancel()
}

func (l *lobbyWorldHandler) HandleBlockBurn(ctx *world.Context, pos cube.Pos) {
	ctx.Cancel()
}

func (l *lobbyWorldHandler) HandleCropTrample(ctx *world.Context, pos cube.Pos) {
	ctx.Cancel()
}

func (l *lobbyWorldHandler) HandleLeavesDecay(ctx *world.Context, pos cube.Pos) {
	ctx.Cancel()
}

func (l *lobbyWorldHandler) HandleEntitySpawn(tx *world.Tx, e world.Entity) {

}

func (l *lobbyWorldHandler) HandleEntityDespawn(tx *world.Tx, e world.Entity) {

}

func (l *lobbyWorldHandler) HandleClose(tx *world.Tx) {

}
