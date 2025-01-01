package practice

import (
	"github.com/akmalfairuz/df-practice/practice/ffa"
	"github.com/akmalfairuz/df-practice/practice/game/igame"
	"github.com/akmalfairuz/df-practice/practice/helper"
	"github.com/akmalfairuz/df-practice/practice/lobby"
	"github.com/akmalfairuz/df-practice/practice/user"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/entity"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/player/skin"
	"github.com/df-mc/dragonfly/server/session"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
	"log/slog"
	"net"
	"time"
)

type playerHandler struct {
	log *slog.Logger
	u   *user.User

	lastChatAt    time.Time
	lastCommandAt time.Time

	l *lobby.Lobby
}

const (
	defaultChatCooldown    = 3 * time.Second
	defaultCommandCooldown = 500 * time.Millisecond

	premiumChatCooldown    = 500 * time.Millisecond
	premiumCommandCooldown = 300 * time.Millisecond
)

func newPlayerHandler(pr *Practice, u *user.User) *playerHandler {
	return &playerHandler{
		u:   u,
		l:   pr.l,
		log: pr.log,
	}
}

func (ph *playerHandler) HandleMove(ctx *player.Context, newPos mgl64.Vec3, newRot cube.Rotation) {
	if ffaArena := ph.ffaArena(ctx.Val()); ffaArena != nil {
		ffaArena.HandleMove(ctx, newPos, newRot)
	}

	if g := ph.game(ctx.Val()); g != nil {
		g.HandleMove(ctx, newPos, newRot)
	}
}

func (ph *playerHandler) HandleJump(p *player.Player) {

}

func (ph *playerHandler) HandleTeleport(ctx *player.Context, pos mgl64.Vec3) {

}

func (ph *playerHandler) HandleChangeWorld(p *player.Player, before, after *world.World) {
	u := user.Get(p)
	if u != nil {
		u.ResetComboCounter()
	}
}

func (ph *playerHandler) HandleToggleSprint(ctx *player.Context, after bool) {

}

func (ph *playerHandler) HandleToggleSneak(ctx *player.Context, after bool) {

}

func (ph *playerHandler) HandleChat(ctx *player.Context, message *string) {
	ctx.Cancel() // Prevent chat handled by dragonfly

	if time.Since(ph.lastChatAt) < defaultChatCooldown {
		ph.u.Messaget("error.cooldown.chat", time.Until(ph.lastChatAt.Add(defaultChatCooldown)).Seconds())
		return
	}
	ph.lastChatAt = time.Now()

	user.BroadcastMessaget("chat.message", ph.u.Name(), *message)
}

func (ph *playerHandler) HandleFoodLoss(ctx *player.Context, from int, to *int) {
	if lobby.Instance().IsInLobby(ctx.Val()) {
		ctx.Cancel()
		return
	}

	if ffaArena := ph.ffaArena(ctx.Val()); ffaArena != nil {
		ffaArena.HandleFoodLoss(ctx, from, to)
	}

	if g := ph.game(ctx.Val()); g != nil {
		g.HandleFoodLoss(ctx, from, to)
	}
}

func (ph *playerHandler) HandleHeal(ctx *player.Context, health *float64, src world.HealingSource) {
	if _, ok := src.(helper.SetHealthSource); ok {
		return
	}

	if ffaArena := ph.ffaArena(ctx.Val()); ffaArena != nil {
		ffaArena.HandleHeal(ctx, health, src)
	}

	if g := ph.game(ctx.Val()); g != nil {
		g.HandleHeal(ctx, health, src)
	}
}

func (ph *playerHandler) HandleHurt(ctx *player.Context, damage *float64, immune bool, attackImmunity *time.Duration, src world.DamageSource) {
	if _, ok := src.(entity.AttackDamageSource); ok {
		*attackImmunity = (time.Millisecond * 50) * 9
	}

	if lobby.Instance().IsInLobby(ctx.Val()) {
		ctx.Cancel()
		return
	}

	if *damage > ctx.Val().Health() {
		ctx.Cancel()
	}

	if ffaArena := ph.ffaArena(ctx.Val()); ffaArena != nil {
		ffaArena.HandleHurt(ctx, damage, immune, attackImmunity, src)
	}

	if g := ph.game(ctx.Val()); g != nil {
		g.HandleHurt(ctx, damage, immune, attackImmunity, src)
	}

	// Final check (*damage may be changed), if death, reset combo counter
	if !immune && *damage > ctx.Val().Health() {
		if u := user.Get(ctx.Val()); u != nil {
			u.ResetComboCounter()
		}
	}

	if attackSource, ok := src.(entity.AttackDamageSource); ok {
		if attacker, ok := attackSource.Attacker.(*player.Player); ok {
			if u := user.Get(attacker); u != nil {
				if !ctx.Cancelled() && !immune {
					u.AddComboCounter()
				}

				// TODO: Basic reach distance, this is not accurate, account for latency and other network factors
				eyePos := attacker.Position().Add(mgl64.Vec3{0, ctx.Val().EyeHeight(), 0})
				targetPos := ctx.Val().Position()
				reachDistance := eyePos.Sub(targetPos).Len()
				u.SetLastHitReachDistance(reachDistance)
			}
			if !immune {
				if u := user.Get(ctx.Val()); u != nil {
					u.ResetComboCounter()
				}
			}
		}
	}
}

func (ph *playerHandler) HandleDeath(p *player.Player, src world.DamageSource, keepInv *bool) {
	panic("HandleDeath: this should never be called")
}

func (ph *playerHandler) HandleRespawn(p *player.Player, pos *mgl64.Vec3, w **world.World) {
	panic("HandleRespawn: this should never be called")
}

func (ph *playerHandler) HandleSkinChange(ctx *player.Context, skin *skin.Skin) {
	user.Get(ctx.Val()).Messaget("error.skin.change.not.allowed")
	ctx.Cancel()
}

func (ph *playerHandler) HandleFireExtinguish(ctx *player.Context, pos cube.Pos) {
	if lobby.Instance().IsInLobby(ctx.Val()) {
		ctx.Cancel()
		return
	}
}

func (ph *playerHandler) HandleStartBreak(ctx *player.Context, pos cube.Pos) {
	if lobby.Instance().IsInLobby(ctx.Val()) {
		ctx.Cancel()
		return
	}
}

func (ph *playerHandler) HandleBlockBreak(ctx *player.Context, pos cube.Pos, drops *[]item.Stack, xp *int) {
	if lobby.Instance().IsInLobby(ctx.Val()) {
		ctx.Cancel()
		return
	}

	if ffaArena := ph.ffaArena(ctx.Val()); ffaArena != nil {
		ffaArena.HandleBlockBreak(ctx, pos, drops, xp)
	}

	if g := ph.game(ctx.Val()); g != nil {
		g.HandleBlockBreak(ctx, pos, drops, xp)
	}
}

func (ph *playerHandler) HandleBlockPlace(ctx *player.Context, pos cube.Pos, b world.Block) {
	if lobby.Instance().IsInLobby(ctx.Val()) {
		ctx.Cancel()
		return
	}

	if ffaArena := ph.ffaArena(ctx.Val()); ffaArena != nil {
		ffaArena.HandleBlockPlace(ctx, pos, b)
	}

	if g := ph.game(ctx.Val()); g != nil {
		g.HandleBlockPlace(ctx, pos, b)
	}
}

func (ph *playerHandler) HandleBlockPick(ctx *player.Context, pos cube.Pos, b world.Block) {
	ctx.Cancel() // Prevent block picking
}

func (ph *playerHandler) HandleItemUse(ctx *player.Context) {
	ph.l.HandleItemUse(ctx)

	if ffaArena := ph.ffaArena(ctx.Val()); ffaArena != nil {
		ffaArena.HandleItemUse(ctx)
	}

	if g := ph.game(ctx.Val()); g != nil {
		g.HandleItemUse(ctx)
	}
}

func (ph *playerHandler) HandleItemUseOnBlock(ctx *player.Context, pos cube.Pos, face cube.Face, clickPos mgl64.Vec3) {
	if lobby.Instance().IsInLobby(ctx.Val()) {
		ctx.Cancel()
		return
	}
}

func (ph *playerHandler) HandleItemUseOnEntity(ctx *player.Context, e world.Entity) {
	if g := ph.game(ctx.Val()); g != nil {
		g.HandleItemUseOnEntity(ctx, e)
	}
}

func (ph *playerHandler) HandleItemRelease(ctx *player.Context, item item.Stack, dur time.Duration) {

}

func (ph *playerHandler) HandleItemConsume(ctx *player.Context, item item.Stack) {
	if lobby.Instance().IsInLobby(ctx.Val()) {
		ctx.Cancel()
		return
	}
}

func (ph *playerHandler) HandleAttackEntity(ctx *player.Context, e world.Entity, force, height *float64, critical *bool) {
	if lobby.Instance().IsInLobby(ctx.Val()) {
		ctx.Cancel()
		return
	}

	// Venity Knockback
	*force = 0.4
	*height = 0.405

	if ffaArena := ph.ffaArena(ctx.Val()); ffaArena != nil {
		if targetPlayer, ok := e.(*player.Player); ok {
			if !ffaArena.IsInArena(user.Get(targetPlayer)) {
				ctx.Cancel()
			}
		}
	}

	if g := ph.game(ctx.Val()); g != nil {
		g.HandleAttackEntity(ctx, e, force, height, critical)
	}
}

func (ph *playerHandler) HandleExperienceGain(ctx *player.Context, amount *int) {
	ctx.Cancel() // Practice server doesn't handle experience gain
}

func (ph *playerHandler) HandlePunchAir(ctx *player.Context) {

}

func (ph *playerHandler) HandleSignEdit(ctx *player.Context, pos cube.Pos, frontSide bool, oldText, newText string) {

}

func (ph *playerHandler) HandleLecternPageTurn(ctx *player.Context, pos cube.Pos, oldPage int, newPage *int) {

}

func (ph *playerHandler) HandleItemDamage(ctx *player.Context, i item.Stack, damage int) {
	if helper.IsItemUnbreakable(i) {
		ctx.Cancel()
	}
}

func (ph *playerHandler) HandleItemPickup(ctx *player.Context, i *item.Stack) {
	if lobby.Instance().IsInLobby(ctx.Val()) {
		ctx.Cancel()
		return
	}

	if g := ph.game(ctx.Val()); g != nil {
		g.HandleItemPickup(ctx, i)
	}
}

func (ph *playerHandler) HandleHeldSlotChange(ctx *player.Context, from, to int) {

}

func (ph *playerHandler) HandleItemDrop(ctx *player.Context, s item.Stack) {
	if lobby.Instance().IsInLobby(ctx.Val()) {
		ctx.Cancel()
		return
	}

	if ffaArena := ph.ffaArena(ctx.Val()); ffaArena != nil {
		ffaArena.HandleItemDrop(ctx, s)
	}
}

func (ph *playerHandler) HandleTransfer(ctx *player.Context, addr *net.UDPAddr) {
	panic("HandleTransfer: this should never be called")
}

func (ph *playerHandler) HandleCommandExecution(ctx *player.Context, command cmd.Command, args []string) {
	if time.Since(ph.lastCommandAt) < defaultCommandCooldown {
		ph.u.Messaget("error.cooldown.command", time.Until(ph.lastCommandAt.Add(defaultCommandCooldown)).Seconds())
		ctx.Cancel()
		return
	}
	ph.lastCommandAt = time.Now()
}

func (ph *playerHandler) HandleQuit(p *player.Player) {
	go ph.u.SynchronizeLastSeen()

	if ffaArena := ph.ffaArena(p); ffaArena != nil {
		helper.LogErrors(ffaArena.Quit(p))
	}

	if g := ph.game(p); g != nil {
		helper.LogErrors(g.Quit(p))
	}

	user.BroadcastMessaget("player.quit.message", p.Name())

	_ = ph.u.Close()
	ph.log.Info("player disconnected", "player", p.Name())
}

func (ph *playerHandler) HandleDiagnostics(p *player.Player, d session.Diagnostics) {

}

func (ph *playerHandler) ffaArena(p *player.Player) *ffa.Arena {
	ret := ph.u.CurrentFFAArena()
	if ret == nil {
		return nil
	}
	return ret.(*ffa.Arena)
}

func (ph *playerHandler) game(p *player.Player) igame.IGame {
	ret := ph.u.CurrentGame()
	if ret == nil {
		return nil
	}
	return ret.(igame.IGame)
}

// Compile-time check to ensure that playerHandler implements player.Handler.
var _ player.Handler = &playerHandler{}
