package schedule

import (
	"fmt"
	"strings"
	"time"

	"StatusApp/configs"
)

// View renders the schedule information.
// It also contains logic for notifications and dynamic styling which could be
// separated further in a more advanced refactoring.
func View(meetings []Meeting) string {
	var sb strings.Builder
	firstLine := true
	numMeetings := 0
	inMeeting := false
	soonMeeting := false
	nowTime, err := time.Parse("15:04", time.Now().Format("15:04"))

	for i, meeting := range meetings {
		if err != nil {
			continue
		}
		if meeting.Time.Compare(nowTime) < 0 {
			if meeting.End.Compare(nowTime) < 0 {
				continue // Meeting is over
			}
			inMeeting = true // Currently in a meeting
			soonMeeting = false
		}

		if firstLine {
			sb.WriteString(
				configs.BoldText.Render(
					fmt.Sprintf(
						"%s - %s  %-58s%s",
						meeting.Time.Format("15:04"),
						meeting.End.Format("15:04"),
						meeting.Title,
						meeting.Room,
					),
				),
			)
			sb.WriteString("\n\n")
			ts := meeting.Time.Sub(nowTime)
			if ts.Minutes() > 0 && ts.Minutes() < 5 {
				soonMeeting = true
			} else {
				soonMeeting = false
			}
			firstLine = false
		} else {
			if i >= 5 {
				break
			}
			fmt.Fprintf(&sb, "%s - %s  %-58s%s", meeting.Time.Format("15:04"), meeting.End.Format("15:04"), meeting.Title, meeting.Room)
			sb.WriteString("\n")
		}
		numMeetings++
	}

	if numMeetings == 0 {
		sb.WriteString("No more meetings today!")
		soonMeeting = false
		inMeeting = false
	}

	// Determine style based on state
	finalStyle := configs.ScheduleStyle
	if inMeeting {
		finalStyle = finalStyle.BorderForeground(configs.HotPink)
	} else if soonMeeting {
		if time.Now().Second()%2 == 0 {
			finalStyle = finalStyle.BorderForeground(configs.StandardText)
		} else {
			finalStyle = finalStyle.BorderForeground(configs.HotPink)
		}
	} else {
		finalStyle = finalStyle.BorderForeground(configs.StandardText)
	}

	return finalStyle.Render(strings.TrimSuffix(sb.String(), "\n"))
}
