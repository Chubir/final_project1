package main

import (
	"strconv"
	"time"
)

func NextDate(now time.Time, date string, repeat string) (string, error) {
	ruleType := repeat[0:1]
	value := repeat[2:]
	switch ruleType {
	case "d":
		days, err := strconv.Atoi(value)
		if err != nil {
			return "", err
		}
		nextDate := now.AddDate(0, 0, days)
		for time.Until(nextDate) < 0 {
			nextDate = nextDate.AddDate(0, 0, days)
		}
		return nextDate.Format(time.DateOnly), nil
	case "y":
		years, _ := strconv.Atoi(value)
		nextDate := now.AddDate(years, 0, 0)
		for time.Until(nextDate) < 0 {
			nextDate = nextDate.AddDate(0, 0, years)
		}
		return nextDate.Format(time.DateOnly), nil
	}
	return date, nil
}
