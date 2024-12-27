package gamemanager

import (
	"github.com/akmalfairuz/df-practice/practice/game/igame"
	"github.com/df-mc/dragonfly/server/player"
	"sync"
)

type Manager struct {
	mu sync.Mutex
	g  map[string]igame.Impl

	createGameFunc func(mgr *Manager) igame.Impl
}

func New(createGameFunc func(mgr *Manager) igame.Impl) *Manager {
	return &Manager{
		g:              make(map[string]igame.Impl),
		createGameFunc: createGameFunc,
	}
}

func (mgr *Manager) PlayersCount() int {
	mgr.mu.Lock()
	defer mgr.mu.Unlock()

	var count int
	for _, d := range mgr.g {
		count += len(d.Game().Participants())
	}
	return count
}

func (mgr *Manager) Join(p *player.Player) error {
	mgr.mu.Lock()
	defer mgr.mu.Unlock()

	for _, d := range mgr.g {
		if err := d.Game().Join(p); err == nil {
			return nil
		}
	}

	d := (mgr.createGameFunc)(mgr)
	mgr.g[d.Game().ID()] = d
	d.Game().SetCloseHook(func() {
		mgr.mu.Lock()
		defer mgr.mu.Unlock()
		delete(mgr.g, d.Game().ID())
	})
	if err := d.Game().Load(); err != nil {
		panic(err)
	}
	return d.Game().Join(p)
}
