package gamemanager

import (
	"github.com/akmalfairuz/df-practice/practice/game/igame"
	"github.com/df-mc/dragonfly/server/player"
	"sync"
)

type Manager struct {
	mu sync.Mutex
	g  map[string]igame.Impl

	createGameFunc func() igame.Impl
}

func New(createGameFunc func() igame.Impl) *Manager {
	return &Manager{
		g:              make(map[string]igame.Impl),
		createGameFunc: createGameFunc,
	}
}

func (mgr *Manager) Join(p *player.Player) error {
	mgr.mu.Lock()
	defer mgr.mu.Unlock()

	for _, d := range mgr.g {
		if err := d.Game().Join(p); err == nil {
			return nil
		}
	}

	d := (mgr.createGameFunc)()
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
