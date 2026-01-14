package renderers

import (
	"StatusApp/configs"
	"StatusApp/internal/schedule"
	"fmt"
	"strings"
	"time"

	"github.com/gen2brain/beeep"
)

func RenderSchedule() string {
	firstLine := true
	numMeetings := 0
	meetings := schedule.LoadSchedule()
	var sb strings.Builder
	configs.InMeeting = false
	for i, meeting := range meetings {
		nowTime, err := time.Parse("15:04", time.Now().Format("15:04"))
		if err != nil {
			continue
		}
		if meeting.Time.Compare(nowTime) < 0 {
			if meeting.End.Compare(nowTime) < 0 {
				continue
			} else {
				configs.InMeeting = true
			}
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
			sb.WriteString("\n")
			ts := meeting.Time.Sub(nowTime)
			if ts.Minutes() < 5 {
				if !configs.SoonMeeting {
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
				configs.SoonMeeting = true
			} else {
				configs.SoonMeeting = false
			}
			firstLine = false
		} else {
			if i >= 5 {
				break
			}
			sb.WriteString(fmt.Sprintf("%s - %s  %-58s%s", meeting.Time.Format("15:04"), meeting.End.Format("15:04"), meeting.Title, meeting.Room))
		}
		sb.WriteString("\n")
		numMeetings = numMeetings + 1
	}
	if numMeetings == 0 {
		sb.WriteString("No more meetings today!")
		configs.SoonMeeting = false
		configs.ScheduleOutlined = false
	}

	if configs.InMeeting {
		configs.ScheduleStyle = configs.ScheduleStyle.BorderForeground(configs.HotPink)
	} else if configs.SoonMeeting {
		if configs.ScheduleOutlined {
			configs.ScheduleStyle = configs.ScheduleStyle.BorderForeground(configs.StandardText)
			configs.ScheduleOutlined = false
		} else {
			configs.ScheduleStyle = configs.ScheduleStyle.BorderForeground(configs.HotPink)
			configs.ScheduleOutlined = true
		}
	} else {
		configs.ScheduleStyle = configs.ScheduleStyle.BorderForeground(configs.StandardText)
	}

	return configs.ScheduleStyle.Render(strings.TrimSuffix(sb.String(), "\n"))
}
