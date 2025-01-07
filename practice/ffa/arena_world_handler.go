package ffa

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
)

type arenaWorldHandler struct{}

func (a arenaWorldHandler) HandleLiquidFlow(ctx *world.Context, from, into cube.Pos, liquid world.Liquid, replaced world.Block) {
	ctx.Cancel()
}

func (a arenaWorldHandler) HandleLiquidDecay(ctx *world.Context, pos cube.Pos, before, after world.Liquid) {
	ctx.Cancel()
}

func (a arenaWorldHandler) HandleLiquidHarden(ctx *world.Context, hardenedPos cube.Pos, liquidHardened, otherLiquid, newBlock world.Block) {
	ctx.Cancel()
}

func (a arenaWorldHandler) HandleSound(ctx *world.Context, s world.Sound, pos mgl64.Vec3) {

}

func (a arenaWorldHandler) HandleFireSpread(ctx *world.Context, from, to cube.Pos) {
	ctx.Cancel()
}

func (a arenaWorldHandler) HandleBlockBurn(ctx *world.Context, pos cube.Pos) {
	ctx.Cancel()
}

func (a arenaWorldHandler) HandleCropTrample(ctx *world.Context, pos cube.Pos) {
	ctx.Cancel()
}

func (a arenaWorldHandler) HandleLeavesDecay(ctx *world.Context, pos cube.Pos) {
	ctx.Cancel()
}

func (a arenaWorldHandler) HandleEntitySpawn(tx *world.Tx, e world.Entity) {

}

func (a arenaWorldHandler) HandleEntityDespawn(tx *world.Tx, e world.Entity) {

}

func (a arenaWorldHandler) HandleClose(tx *world.Tx) {

}
