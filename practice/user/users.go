package user

import (
	"github.com/akmalfairuz/df-practice/practice/database"
	"github.com/akmalfairuz/df-practice/practice/repository"
	"github.com/df-mc/dragonfly/server/player"
	"sync"
)

var (
	users          = make(map[string]*User)
	usersMu        sync.RWMutex
	userRepository = repository.NewUser(database.Get())
)

func Get(p *player.Player) *User {
	usersMu.RLock()
	defer usersMu.RUnlock()
	ret, ok := users[p.XUID()]
	if !ok {
		return nil
	}
	return ret
}

func GetByXUID(xuid string) *User {
	usersMu.RLock()
	defer usersMu.RUnlock()
	ret, ok := users[xuid]
	if !ok {
		return nil
	}
	return ret
}

func Remove(p *player.Player) {
	usersMu.Lock()
	defer usersMu.Unlock()
	delete(users, p.XUID())
}

func RemoveByXUID(xuid string) {
	usersMu.Lock()
	defer usersMu.Unlock()
	delete(users, xuid)
}

func Store(u *User) {
	usersMu.Lock()
	defer usersMu.Unlock()
	users[u.xuid] = u
}

func Count() int {
	usersMu.RLock()
	defer usersMu.RUnlock()
	return len(users)
}

func BroadcastMessaget(translationName string, args ...any) {
	usersMu.RLock()
	defer usersMu.RUnlock()
	for _, u := range users {
		u.Messaget(translationName, args...)
	}
}

func BulkMessaget(users []*User, translationName string, args ...any) {
	for _, u := range users {
		u.Messaget(translationName, args...)
	}
}
