package entity

import (
	"github.com/df-mc/dragonfly/server/block/cube/trace"
	"github.com/df-mc/dragonfly/server/entity"
	"github.com/df-mc/dragonfly/server/event"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/particle"
	"github.com/df-mc/dragonfly/server/world/sound"
	"github.com/go-gl/mathgl/mgl64"
)

type PearlTeleportHandler interface {
	HandlePearlTeleport(ctx *player.Context, target mgl64.Vec3)
}

func NewNoDamageEnderPearl(opts world.EntitySpawnOpts, owner world.Entity) *world.EntityHandle {
	conf := entity.ProjectileBehaviourConfig{
		Gravity:  0.03,
		Drag:     0.01,
		Particle: particle.EndermanTeleport{},
		Sound:    sound.Teleport{},
		Hit: func(e *entity.Ent, tx *world.Tx, target trace.Result) {
			owner, _ := e.Behaviour().(*entity.ProjectileBehaviour).Owner().Entity(tx)
			if user, ok := owner.(*player.Player); ok {
				ctx := event.C(user)
				if h, ok := user.Handler().(PearlTeleportHandler); ok {
					h.HandlePearlTeleport(ctx, target.Position())
				}

				if !ctx.Cancelled() {
					tx.PlaySound(user.Position(), sound.Teleport{})
					user.Teleport(target.Position())
				}
			}
		},
	}
	conf.Owner = owner.H()
	return opts.New(entity.EnderPearlType, conf)
}
