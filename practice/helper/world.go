package helper

import (
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	_ "unsafe"
)

func InvalidPlayerCtxWorld(ctx *player.Context, w *world.World) bool {
	if ctx.Val().Tx().World() == w {
		return false
	}
	ctx.Cancel()
	return true
}
