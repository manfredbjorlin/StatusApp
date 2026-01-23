package weather

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// parseYrIcon manually parses a unicode surrogate pair string like `\udb81\udd99`.
// The standard strconv.Unquote is not working reliably with this format.
func parseYrIcon(raw string) (string, error) {
	s := strings.TrimSpace(raw)
	if len(s) != 12 || !strings.HasPrefix(s, `\u`) || s[6:8] != `\u` {
		return "", fmt.Errorf("invalid icon format: expected '\\uXXXX\\uYYYY', got '%s'", raw)
	}

	highStr := s[2:6]
	lowStr := s[8:12]

	high, err := strconv.ParseInt(highStr, 16, 32)
	if err != nil {
		return "", fmt.Errorf("cannot parse high surrogate '%s': %w", highStr, err)
	}

	low, err := strconv.ParseInt(lowStr, 16, 32)
	if err != nil {
		return "", fmt.Errorf("cannot parse low surrogate '%s': %w", lowStr, err)
	}

	// Combine the surrogates to get the final code point.
	// Formula from Unicode standard.
	r := 0x10000 + (high-0xD800)*0x400 + (low-0xDC00)

	return string(rune(r)), nil
}

func getIconYr(csvLocation string, symbolId string) (string, string) {
	file, err := os.Open(csvLocation)
	if err != nil {
		return "?", "Unknown"
	}
	defer file.Close()

	reader := csv.NewReader(file)
	// skip header
	if _, err := reader.Read(); err != nil {
		return "?", "Unknown"
	}

	records, err := reader.ReadAll()
	if err != nil {
		return "?", "Unknown"
	}

	for _, record := range records {
		if strings.TrimSpace(record[0]) == symbolId {
			icon, err := parseYrIcon(record[6]) // Use the new manual parser
			if err != nil {
				return "?", err.Error() // Return the detailed error
			}
			return icon, record[1]
		}
	}

	return "?", "Unknown"
}