package command

import (
	"fmt"
	"github.com/akmalfairuz/df-practice/practice/lang"
	"github.com/akmalfairuz/df-practice/practice/user"
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/go-gl/mathgl/mgl64"
	"golang.org/x/text/language"
)

func translatef(s cmd.Source, key string, args ...interface{}) string {
	if p, ok := s.(*player.Player); ok {
		u := user.Get(p)
		if u != nil {
			return u.Translatef(key, args...)
		}
	}
	return lang.Translatef(language.English, key, args...)
}

func messaget(s cmd.Source, key string, args ...interface{}) {
	s2, ok := s.(interface {
		cmd.Source

		Message(a ...any)
	})

	if !ok {
		return
	}

	s2.Message(translatef(s, key, args...))
}

func isPlayer(s cmd.Source) bool {
	_, ok := s.(*player.Player)
	return ok
}

type onlyPlayer struct{}

func (onlyPlayer) Allow(s cmd.Source) bool {
	return isPlayer(s)
}

type onlyAdmin struct{}

func (onlyAdmin) Allow(s cmd.Source) bool {
	if p, ok := s.(*player.Player); ok {
		u := user.Get(p)
		if u == nil {
			return false
		}
		return u.RankName() == "admin"
	}
	return false
}

func formatPos(pos mgl64.Vec3) string {
	return fmt.Sprintf("x=%.2f, y=%.2f, z=%.2f", pos.X(), pos.Y(), pos.Z())
}
