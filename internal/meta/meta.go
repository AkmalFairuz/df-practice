package meta

import "sync"

var mu sync.RWMutex
var meta = map[string]any{}

func Set(key string, value any) {
	mu.Lock()
	meta[key] = value
	mu.Unlock()
}

func Get(key string) any {
	mu.RLock()
	defer mu.RUnlock()
	return meta[key]
}
