package command

import (
	"github.com/akmalfairuz/df-practice/practice/user"
	"github.com/akmalfairuz/df-practice/translations"
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
)

type Reply struct {
	onlyPlayer

	Message cmd.Varargs `cmd:"message"`
}

func (r Reply) Run(s cmd.Source, o *cmd.Output, tx *world.Tx) {
	u := user.Get(s.(*player.Player))

	replyToXUID := u.ReplyWhisperToXUID()
	if replyToXUID == "" {
		o.Error(translatef(s, translations.ErrorCommandReplyNoTarget))
		return
	}

	target := user.GetByXUID(replyToXUID)
	if target == nil {
		o.Error(translatef(s, translations.ErrorCommandReplyTargetOffline))
		return
	}

	u.OnSendWhisper(target, string(r.Message))
	target.OnReceiveWhisper(u, string(r.Message))
}
