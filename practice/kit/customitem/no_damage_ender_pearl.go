package customitem

import (
	"github.com/akmalfairuz/df-practice/practice/entity"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"
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

	create := entity.NewNoDamageEnderPearl
	opts := world.EntitySpawnOpts{Position: user.Position().Add(mgl64.Vec3{0, eyeHeight}), Velocity: user.Rotation().Vec3().Mul(3)}
	tx.AddEntity(create(opts, user))
	tx.PlaySound(user.Position(), sound.ItemThrow{})

	ctx.SubtractFromCount(1)
	return true
}
