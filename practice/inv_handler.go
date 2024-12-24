package practice

import (
	"github.com/akmalfairuz/df-practice/practice/lobby"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/item/inventory"
	"github.com/df-mc/dragonfly/server/player"
)

type playerInvHandler struct {
	inv *inventory.Inventory
}

func newPlayerInvHandler(inv *inventory.Inventory) *playerInvHandler {
	return &playerInvHandler{
		inv: inv,
	}
}

// Compile-time check to ensure the type implements the Handler interface.
var _ inventory.Handler = &playerInvHandler{}

func (p *playerInvHandler) player(ctx *inventory.Context) *player.Player {
	return ctx.Val().(*player.Player)
}

func (p *playerInvHandler) HandleTake(ctx *inventory.Context, slot int, stack item.Stack) {
	if lobby.Instance().IsInLobby(p.player(ctx)) {
		ctx.Cancel()
	}
}

func (p *playerInvHandler) HandlePlace(ctx *inventory.Context, slot int, stack item.Stack) {
	if lobby.Instance().IsInLobby(p.player(ctx)) {
		ctx.Cancel()
	}
}

func (p *playerInvHandler) HandleDrop(ctx *inventory.Context, slot int, stack item.Stack) {
	if lobby.Instance().IsInLobby(p.player(ctx)) {
		ctx.Cancel()
	}
}
