package main

import (
	"fmt"
	"time"
)

func FormatTime(currentTime *time.Time, layout string) string {
	var timeString string

	switch layout {
	case "Seconds1970":
		timeString = fmt.Sprintf("%d", currentTime.Unix())
	case "Seconds1900":
		timeString = fmt.Sprintf("%d", currentTime.Unix()+2208988800)
	default:
		timeString = currentTime.Format(layout)
	}

	return timeString
}
