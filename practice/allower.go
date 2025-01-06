package practice

import (
	"github.com/akmalfairuz/df-practice/practice/helper"
	"github.com/akmalfairuz/df-practice/practice/lang"
	"github.com/akmalfairuz/df-practice/practice/repository"
	"github.com/akmalfairuz/df-practice/translations"
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
		return lang.Translatef(l, translations.ErrorUnknown), false
	}

	b, err := repository.BanRepo().FindByPlayerID(u.ID)
	if err != nil {
		if repository.IsNotExists(err) {
			return "", true
		}
		helper.LogErrors(err)
		return lang.Translatef(l, translations.ErrorUnknown), false
	}

	now := time.Now().Unix()
	if b.ExpireAt > now {
		daysRemaining, hoursRemaining, minutesRemaining := b.Remaining()
		return strutils.CenterLine(lang.Translatef(l, translations.BannedKickMessage, b.Reason, lang.Translatef(l, translations.TimeShortDhm, daysRemaining, hoursRemaining, minutesRemaining))), false
	}

	helper.LogErrors(repository.BanRepo().DeleteByPlayerID(u.ID))

	return "", true
}
