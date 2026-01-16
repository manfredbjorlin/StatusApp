package models

import "time"

type Meeting struct {
	Time  time.Time
	End   time.Time
	Title string
	Room  string
	Rooms []string
}
