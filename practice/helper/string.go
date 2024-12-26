package helper

import "fmt"

func FormatTime(seconds int) string {
	minutes := seconds / 60
	seconds %= 60
	return fmt.Sprintf("%02d:%02d", minutes, seconds)
}
