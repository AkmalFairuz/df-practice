package command

import (
	"errors"
	"github.com/akmalfairuz/df-practice/practice/user"
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/go-gl/mathgl/mgl64"
	"reflect"
)

type onlineTarget struct {
	name string
	u    *user.User
}

func (o onlineTarget) Type() string {
	return "target"
}

func (o onlineTarget) User() *user.User {
	return o.u
}

func (o onlineTarget) Position() mgl64.Vec3 {
	return mgl64.Vec3{}
}

func (o onlineTarget) ExecutePlayer(f func(p *player.Player, ok bool)) {
	o.u.ExecutePlayer(func(p *player.Player, ok bool) {
		f(p, ok)
	})
}

func (o onlineTarget) Parse(line *cmd.Line, v reflect.Value) error {
	next, ok := line.Next()
	if ok {
		u, ok := user.GetByPrefix(next)
		if ok {
			v.Set(reflect.ValueOf(onlineTarget{
				name: next,
				u:    u,
			}))
			return nil
		}
		return errors.New("target not found")
	}
	return errors.New("no target provided")
}
