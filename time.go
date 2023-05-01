package main

import (
	"fmt"
	"time"
)

func parseTime(timeStr string) (string, error) {
	formats := []string{
		"Mon, 02 Jan 2006 15:04:05 -0700",
		"2006-01-02T15:04:05.999-07:00",
		"Mon, 02 Jan 2006 15:04:05 +0700",
		"Mon, 02 Jan 2006 15:04:05 GMT",
		"Mon, 02 Jan 2006 15:04:05 -0700",
	}

	for _, format := range formats {
		t, err := time.Parse(format, timeStr)
		if err == nil {
			return t.UTC().Format(time.RFC3339), nil
		}
	}

	return "", fmt.Errorf("Could not parse time: %s", timeStr)
}
