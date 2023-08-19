package pkg

import (
	"fmt"
	"time"
)

func StartOfMonth() time.Time {
	now := time.Now()
	return time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
}

func ParseDate(date string) (time.Time, error) {
	now := time.Now()
	now = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	switch date {
	case "", "today", "Today":
		return now, nil
	case "yesterday", "Yesterday":
		return now.AddDate(0, 0, -1), nil
	case "tomorrow", "Tomorrow":
		return now.AddDate(0, 0, 1), nil
	}

	dateFormats := []string{time.DateOnly, "02-01-2006", "Jan 2, 2006", "January 2, 2006"}
	for _, format := range dateFormats {
		t, err := time.ParseInLocation(format, date, now.Location())
		if err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("invalid date format")
}

func ParseTime(tim string) (time.Time, error) {
	now := time.Now()
	switch tim {
	case "", "now", "Now":
		return now, nil
	case "midnight", "Midnight":
		return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()), nil
	case "morning", "Morning":
		return time.Date(now.Year(), now.Month(), now.Day(), 6, 0, 0, 0, now.Location()), nil
	case "noon", "Noon":
		return time.Date(now.Year(), now.Month(), now.Day(), 12, 0, 0, 0, now.Location()), nil
	case "afternoon", "Afternoon":
		return time.Date(now.Year(), now.Month(), now.Day(), 15, 0, 0, 0, now.Location()), nil
	case "evening", "Evening":
		return time.Date(now.Year(), now.Month(), now.Day(), 18, 0, 0, 0, now.Location()), nil
	case "night", "Night":
		return time.Date(now.Year(), now.Month(), now.Day(), 22, 0, 0, 0, now.Location()), nil
	}

	timeFormats := []string{time.TimeOnly, time.Kitchen, "3:04pm", "3:04 PM", "3:04 pm", "3:04", "15:04"}
	for _, format := range timeFormats {
		t, err := time.ParseInLocation(format, tim, now.Location())
		if err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("invalid time format")
}
