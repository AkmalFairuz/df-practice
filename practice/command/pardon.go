package command

import (
	"github.com/akmalfairuz/df-practice/practice/helper"
	"github.com/akmalfairuz/df-practice/practice/repository"
	"github.com/akmalfairuz/df-practice/translations"
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/world"
)

type Pardon struct {
	onlyAdmin
	Target OfflineTarget `cmd:"target"`
}

func (p Pardon) Run(s cmd.Source, o *cmd.Output, tx *world.Tx) {
	go func() {
		u, err := repository.UserRepo().FindByName(string(p.Target))
		if err != nil {
			if repository.IsNotExists(err) {
				messaget(s, translations.ErrorCommandPardonTargetNotFound)
				return
			}
			helper.LogErrors(err)
			messaget(s, translations.ErrorUnknown)
			return
		}

		rowsAffected, err := repository.BanRepo().DeleteByPlayerID(u.ID)
		if err != nil {
			helper.LogErrors(err)
			messaget(s, translations.ErrorUnknown)
			return
		}

		if rowsAffected == 0 {
			messaget(s, translations.ErrorCommandPardonTargetNotBanned)
			return
		}

		messaget(s, translations.CommandPardonSuccess, u.Name)
	}()
}
