package schedule

import (
	"os"
	"testing"
	"time"
)

func TestLoadSchedule(t *testing.T) {
	content := []byte(`10:00##Test Meeting 1##Room 1##11:00
14:00##Test Meeting 2##Teams##15:00
09:00##Meeting That Should Be First##Room 2##09:30
`)
	tmpfile, err := os.CreateTemp("", "schedule-*.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write(content); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	meetings, err := LoadSchedule(tmpfile.Name())
	if err != nil {
		t.Fatalf("LoadSchedule failed: %v", err)
	}

	if len(meetings) != 3 {
		t.Fatalf("expected 3 meetings, got %d", len(meetings))
	}

	// Test sorting
	expectedFirstTitle := "Meeting That Should Be First"
	if meetings[0].Title != expectedFirstTitle {
		t.Errorf("expected first meeting to be '%s', got '%s'", expectedFirstTitle, meetings[0].Title)
	}

	// Test parsing on a specific meeting after sorting
	var foundMeeting Meeting
	for _, m := range meetings {
		if m.Title == "Test Meeting 2" {
			foundMeeting = m
			break
		}
	}

	if foundMeeting.Title == "" {
		t.Fatal("could not find 'Test Meeting 2' in the loaded meetings")
	}

	expectedTime, _ := time.Parse("15:04", "14:00")
	if !foundMeeting.Time.Equal(expectedTime) {
		t.Errorf("expected meeting time to be 14:00, got %s", foundMeeting.Time.Format("15:04"))
	}
	if foundMeeting.Room != "Teams" {
		t.Errorf("expected room 'Teams', got '%s'", foundMeeting.Room)
	}
}

func TestLoadSchedule_fileNotFound(t *testing.T) {
	_, err := LoadSchedule("non-existent-file.txt")
	if err == nil {
		t.Fatal("expected an error for a missing schedule file, but got none")
	}
}
