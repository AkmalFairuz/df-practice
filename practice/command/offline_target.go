package command

import (
	"errors"
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/go-gl/mathgl/mgl64"
	"reflect"
)

type OfflineTarget string

// Position ...
func (o OfflineTarget) Position() mgl64.Vec3 {
	return mgl64.Vec3{}
}

func (o OfflineTarget) Parse(line *cmd.Line, v reflect.Value) error {
	str, ok := line.Next()
	if !ok {
		return errors.New("not enough arguments")
	}
	v.SetString(str)
	return nil
}

func (o OfflineTarget) Type() string {
	return "target"
}
