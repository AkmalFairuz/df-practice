package user

import (
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/session"
	"reflect"
	"unsafe"
	_ "unsafe"
)

// noinspection ALL
//
//go:linkname playerSession github.com/df-mc/dragonfly/server/player.(*Player).session
func playerSession(*player.Player) *session.Session

func sessionConn(s *session.Session) session.Conn {
	rf := reflect.ValueOf(s).Elem().FieldByName("conn")
	return reflect.NewAt(rf.Type(), unsafe.Pointer(rf.UnsafeAddr())).Elem().Interface().(session.Conn)
}
