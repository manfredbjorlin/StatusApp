package schedule

import "time"

// Meeting represents a single scheduled event.
type Meeting struct {
	Time  time.Time
	End   time.Time
	Title string
	Room  string
	Rooms []string
}
