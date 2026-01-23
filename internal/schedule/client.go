package schedule

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"
)

func LoadSchedule(filePath string) ([]Meeting, error) {
	m := make([]Meeting, 0)

	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("could not find or open schedule file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		meeting := Meeting{}
		if len(scanner.Text()) < 2 {
			continue
		}

		meetingParts := strings.Split(scanner.Text(), "##")
		if len(meetingParts) < 4 {
			continue // Skip malformed lines
		}

		mt := meetingParts[0]
		meetingTime, err := time.Parse("15:04", mt)
		if err != nil {
			continue
		}
		meeting.Time = meetingTime
		meeting.End, _ = time.Parse("15:04", meetingParts[3])
		meeting.Title = meetingParts[1]

		rooms := strings.Split(meetingParts[2], ";")
		currentRoom := strings.TrimSpace(meetingParts[2])

		if len(rooms) > 1 {
			for _, room := range rooms {
				if strings.HasPrefix(room, "Microsoft Teams") {
					continue
				}
				currentRoom = strings.TrimSpace(room)
			}
		}

		if strings.Contains(currentRoom, "M OSL Schweigaards") {
			p := strings.Split(currentRoom, " ")
			if len(p) > 5 {
				currentRoom = p[5]
			}
		} else if strings.Contains(currentRoom, "Microsoft Teams") {
			currentRoom = "Teams"
		}
		meeting.Room = currentRoom
		meeting.Rooms = rooms

		m = append(m, meeting)
	}

	sort.Slice(m, func(i, j int) bool {
		return m[i].Time.Before(m[j].Time)
	})

	return m, nil
}
