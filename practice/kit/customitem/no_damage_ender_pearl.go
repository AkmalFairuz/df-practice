package customitem

import (
	"github.com/df-mc/dragonfly/server/block/cube/trace"
	"github.com/df-mc/dragonfly/server/entity"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/particle"
	"github.com/df-mc/dragonfly/server/world/sound"
	"github.com/go-gl/mathgl/mgl64"
)

type NoDamageEnderPearl struct {
	item.EnderPearl
}

func (e NoDamageEnderPearl) Use(tx *world.Tx, user item.User, ctx *item.UseContext) bool {
	eyeHeight := float64(0)
	if user.(interface{ EyeHeight() float64 }) != nil {
		eyeHeight = user.(interface{ EyeHeight() float64 }).EyeHeight()
	}

	create := newNoDamageEnderPearl
	opts := world.EntitySpawnOpts{Position: user.Position().Add(mgl64.Vec3{0, eyeHeight}), Velocity: user.Rotation().Vec3().Mul(2.2)}
	tx.AddEntity(create(opts, user))
	tx.PlaySound(user.Position(), sound.ItemThrow{})

	ctx.SubtractFromCount(1)
	return true
}

func newNoDamageEnderPearl(opts world.EntitySpawnOpts, owner world.Entity) *world.EntityHandle {
	conf := entity.ProjectileBehaviourConfig{
		Gravity:  0.03,
		Drag:     0.01,
		Particle: particle.EndermanTeleport{},
		Sound:    sound.Teleport{},
		Hit: func(e *entity.Ent, tx *world.Tx, target trace.Result) {
			owner, _ := e.Behaviour().(*entity.ProjectileBehaviour).Owner().Entity(tx)
			if user, ok := owner.(*player.Player); ok {
				tx.PlaySound(user.Position(), sound.Teleport{})
				user.Teleport(target.Position())
			}
		},
	}
	conf.Owner = owner.H()
	return opts.New(entity.EnderPearlType, conf)
}
