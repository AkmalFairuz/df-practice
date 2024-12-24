package practice

import (
	"github.com/akmalfairuz/df-practice/practice/user"
	"time"
)

func startPlayerTick(u *user.User) {
	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()

	currentTick := int64(0)

	for {
		select {
		case <-ticker.C:
			if u.Closed() {
				return
			}
			currentTick++
			handlePlayerTick(u, currentTick)
		}
	}
}

func handlePlayerTick(u *user.User, currentTick int64) {
}
