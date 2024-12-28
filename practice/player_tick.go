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
	if currentTick%20 == 0 {
		//u.EntityHandle().ExecWorld(func(tx *world.Tx, e world.Entity) {
		//	e.(*player.Player).SendTip(fmt.Sprintf("X: %.1f Y: %.1f Z: %.1f YAW: %.0f PIT: %.0f", e.Position().X(), e.Position().Y(), e.Position().Z(), e.Rotation().Yaw(), e.Rotation().Pitch()))
		//})
		u.RemoveOldClicks()

		if time.Since(u.LastComboCounterModified()) > 5*time.Second {
			u.ResetComboCounter()
		}
	}
}
