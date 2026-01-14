package schedule

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"StatusApp/internal/models"
)

func LoadSchedule() []models.Meeting {
	m := make([]models.Meeting, 0)

	fileLocation := os.Getenv("SCHEDULE_FILE_PATH")
	file, err := os.Open(fileLocation)
	if err != nil {
		panic("Could not find schedule file")
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic(fmt.Sprintf("error closing file: %v", err))
		}
	}(file)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		meeting := models.Meeting{}
		if len(scanner.Text()) < 2 {
			continue
		}

		meetingParts := strings.Split(scanner.Text(), "##")
		mt := meetingParts[0]
		meetingTime, err := time.Parse("15:04", mt)
		if err != nil {
			continue
		}
		meeting.Time = meetingTime
		meeting.End, _ = time.Parse("15:04", meetingParts[3])
		meeting.Title = meetingParts[1]
		if len(meeting.Title) > 55 {
			meeting.Title = meeting.Title[:50] + "..."
		}

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
			currentRoom = p[5]
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

	return m
}
