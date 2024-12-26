package practice

import (
	"github.com/akmalfairuz/df-practice/practice/game"
	"github.com/akmalfairuz/df-practice/practice/lobby"
	"github.com/akmalfairuz/df-practice/practice/user"
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

	if g := p.game(ctx); g != nil {
		g.HandleTake(ctx, slot, stack)
	}
}

func (p *playerInvHandler) HandlePlace(ctx *inventory.Context, slot int, stack item.Stack) {
	if lobby.Instance().IsInLobby(p.player(ctx)) {
		ctx.Cancel()
	}

	if g := p.game(ctx); g != nil {
		g.HandlePlace(ctx, slot, stack)
	}
}

func (p *playerInvHandler) HandleDrop(ctx *inventory.Context, slot int, stack item.Stack) {
	if lobby.Instance().IsInLobby(p.player(ctx)) {
		ctx.Cancel()
	}

	if g := p.game(ctx); g != nil {
		g.HandleDrop(ctx, slot, stack)
	}
}

func (p *playerInvHandler) game(ctx *inventory.Context) *game.Game {
	g := user.Get(p.player(ctx)).CurrentGame()
	if g == nil {
		return nil
	}
	return g.(*game.Game)
}
