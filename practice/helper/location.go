package helper

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/go-gl/mathgl/mgl64"
	"reflect"
	"unsafe"
)

type Location struct {
	X     float64 `yaml:"x"`
	Y     float64 `yaml:"y"`
	Z     float64 `yaml:"z"`
	Yaw   float32 `yaml:"yaw"`
	Pitch float32 `yaml:"pitch"`
}

func (loc Location) ToMgl64Vec3() mgl64.Vec3 {
	return mgl64.Vec3{loc.X, loc.Y, loc.Z}
}

func (loc Location) ToMgl32Vec3() mgl32.Vec3 {
	return mgl32.Vec3{float32(loc.X), float32(loc.Y), float32(loc.Z)}
}

func (loc Location) TeleportPlayer(p *player.Player) {
	// TODO: don't use reflect when dragonfly has a method to set player rotation
	rf := reflect.ValueOf(p).Elem().FieldByName("data")
	reflect.NewAt(rf.Type(), unsafe.Pointer(rf.UnsafeAddr())).Elem().Elem().FieldByName("Rot").Set(reflect.ValueOf(cube.Rotation{float64(loc.Yaw), float64(loc.Pitch)}))

	p.Teleport(loc.ToMgl64Vec3())
}

func Mgl64Vec3ToMgl32Vec3(v mgl64.Vec3) mgl32.Vec3 {
	return mgl32.Vec3{float32(v.X()), float32(v.Y()), float32(v.Z())}
}
