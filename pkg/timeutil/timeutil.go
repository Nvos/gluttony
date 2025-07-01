package timeutil

import (
	"fmt"
	"time"
)

func FormatDuration(d time.Duration) string {
	const hourMinutes = 60

	hours := int(d.Hours())
	minutes := int(d.Minutes()) % hourMinutes
	return fmt.Sprintf("%02d:%02d", hours, minutes)
}
