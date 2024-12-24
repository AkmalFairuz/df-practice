package command

import (
	"github.com/akmalfairuz/df-practice/practice/user"
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
)

type Whisper struct {
	onlyPlayer

	Target  []cmd.Target `cmd:"target"`
	Message cmd.Varargs  `cmd:"message"`
}

func (w Whisper) Run(s cmd.Source, o *cmd.Output, tx *world.Tx) {
	u := user.Get(s.(*player.Player))

	if len(w.Target) == 0 {
		o.Error(translatef(s, "error.command.whisper.missing.target"))
		return
	}

	if len(w.Message) == 0 {
		o.Error(translatef(s, "error.command.whisper.missing.message"))
		return
	}

	target := user.Get(w.Target[0].(*player.Player))
	u.OnSendWhisper(target, string(w.Message))
	target.OnReceiveWhisper(u, string(w.Message))
}
