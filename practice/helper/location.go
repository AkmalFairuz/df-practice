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
	X     float64
	Y     float64
	Z     float64
	Yaw   float32
	Pitch float32
}

func ParseSliceOfLocation(slice [][5]float64) []Location {
	var locs []Location
	for _, loc := range slice {
		locs = append(locs, Location{loc[0], loc[1], loc[2], float32(loc[3]), float32(loc[4])})
	}
	return locs
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
