package practice

import (
	"github.com/akmalfairuz/df-practice/practice/helper"
	"github.com/akmalfairuz/df-practice/practice/lang"
	"github.com/akmalfairuz/df-practice/practice/repository"
	"github.com/sandertv/gophertunnel/minecraft/protocol/login"
	"github.com/unickorn/strutils"
	"net"
	"time"
)

type Allower struct{}

func NewAllower() *Allower {
	return &Allower{}
}

func (a *Allower) Allow(_ net.Addr, d login.IdentityData, c login.ClientData) (string, bool) {
	l := lang.ToLangTag(c.LanguageCode)

	u, err := repository.UserRepo().FindByXUID(d.XUID)
	if err != nil {
		if repository.IsNotExists(err) {
			return "", true
		}
		helper.LogErrors(err)
		return lang.Translatef(l, "error.unknown"), false
	}

	b, err := repository.BanRepo().FindByPlayerID(u.ID)
	if err != nil {
		if repository.IsNotExists(err) {
			return "", true
		}
		helper.LogErrors(err)
		return lang.Translatef(l, "error.unknown"), false
	}

	now := time.Now().Unix()
	if b.ExpireAt > now {
		daysRemaining, hoursRemaining, minutesRemaining := b.Remaining()
		return strutils.CenterLine(lang.Translatef(l, "banned.kick.message", b.Reason, lang.Translatef(l, "time.short.dhm", daysRemaining, hoursRemaining, minutesRemaining))), false
	}

	helper.LogErrors(repository.BanRepo().DeleteByPlayerID(u.ID))

	return "", true
}
