package schedule

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/gen2brain/beeep"

	"StatusApp/configs"
)

// View renders the schedule information.
// It also contains logic for notifications and dynamic styling which could be
// separated further in a more advanced refactoring.
func View(meetings []Meeting, inMeeting *bool, soonMeeting *bool, outlined *bool) string {
	var sb strings.Builder
	firstLine := true
	numMeetings := 0
	*inMeeting = false

	for i, meeting := range meetings {
		nowTime, err := time.Parse("15:04", time.Now().Format("15:04"))
		if err != nil {
			continue
		}
		if meeting.Time.Compare(nowTime) < 0 {
			if meeting.End.Compare(nowTime) < 0 {
				continue // Meeting is over
			}
			*inMeeting = true // Currently in a meeting
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
				if !*soonMeeting {
					beeep.AppName = "Meeting Notifier"
					_ = beeep.Notify(
						fmt.Sprintf(
							"%s - %s - %s",
							meeting.Time.Format("15:04"),
							meeting.End.Format("15:04"),
							meeting.Title,
						),
						"",
						nil,
					)
				}
				*soonMeeting = true
			} else {
				*soonMeeting = false
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
		*soonMeeting = false
		*outlined = false
	}

	// Determine style based on state
	var finalStyle lipgloss.Style
	if *inMeeting {
		finalStyle = configs.ScheduleStyle.Copy().BorderForeground(configs.HotPink)
	} else if *soonMeeting {
		if *outlined {
			finalStyle = configs.ScheduleStyle.Copy().BorderForeground(configs.StandardText)
			*outlined = false
		} else {
			finalStyle = configs.ScheduleStyle.Copy().BorderForeground(configs.HotPink)
			*outlined = true
		}
	} else {
		finalStyle = configs.ScheduleStyle.Copy().BorderForeground(configs.StandardText)
	}

	return finalStyle.Render(strings.TrimSuffix(sb.String(), "\n"))
}
