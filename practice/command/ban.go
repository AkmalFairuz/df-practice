package command

import (
	"github.com/akmalfairuz/df-practice/practice/helper"
	"github.com/akmalfairuz/df-practice/practice/model"
	"github.com/akmalfairuz/df-practice/practice/repository"
	"github.com/akmalfairuz/df-practice/practice/user"
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/unickorn/strutils"
	"time"
)

type Ban struct {
	onlyAdmin

	Target   OfflineTarget `cmd:"target"`
	Duration int           `cmd:"duration"`
	Reason   cmd.Varargs   `cmd:"reason"`
}

func (b Ban) Run(s cmd.Source, o *cmd.Output, tx *world.Tx) {
	go func() {
		u, err := repository.UserRepo().FindByName(string(b.Target))
		if err != nil {
			if repository.IsNotExists(err) {
				messaget(s, "error.command.ban.target.not.found")
				return
			}
			helper.LogErrors(err)
			messaget(s, "error.unknown")
			return
		}

		currentBan, err := repository.BanRepo().FindByPlayerID(u.ID)
		if err != nil && !repository.IsNotExists(err) {
			helper.LogErrors(err)
			messaget(s, "error.unknown")
			return
		}

		if currentBan.ID != 0 {
			messaget(s, "error.command.ban.target.already.banned")
			return
		}

		expiresAt := time.Now().Add(time.Duration(b.Duration) * time.Hour * 24)
		ba := model.Ban{
			PlayerID:  u.ID,
			Reason:    string(b.Reason),
			ExpireAt:  expiresAt.Unix(),
			CreatedAt: time.Now().Unix(),
		}
		if _, err := repository.BanRepo().Create(ba); err != nil {
			helper.LogErrors(err)
			messaget(s, "error.unknown")
			return
		}

		messaget(s, "command.ban.success", u.Name, b.Reason, expiresAt.Format("2006-01-02 15:04:05"))

		t := user.GetByXUID(u.XUID)
		if t != nil {
			daysRemaining, hoursRemaining, minutesRemaining := ba.Remaining()
			t.Disconnect(strutils.CenterLine(t.Translatef("banned.kick.message", b.Reason, t.Translatef("time.short.dhm", daysRemaining, hoursRemaining, minutesRemaining))))
		}
	}()
}
