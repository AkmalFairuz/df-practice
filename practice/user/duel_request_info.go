package user

import "time"

type DuelRequestInfo struct {
	TargetXUID string
	RequestAt  time.Time
}
