package common

import (
	"fmt"
	"time"
)

func GetTimeDifferenceString(myTime time.Time) string {
	offlineDiff := time.Since(myTime)
	if time.Now().Before(myTime) {
		offlineDiff = time.Until(myTime)
	}
	diffText := ""
	if offlineDiff.Hours() >= 24 {
		days := int(offlineDiff.Hours()) / 24
		diffText = fmt.Sprintf("%d d", days)
	} else if offlineDiff.Hours() >= 1 {
		hours := int(offlineDiff.Hours())
		diffText = fmt.Sprintf("%d H", hours)
	} else {
		minutes := int(offlineDiff.Minutes())
		diffText = fmt.Sprintf("%d m", minutes)
	}

	return diffText
}
