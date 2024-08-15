package main

import (
	"strconv"
	"time"
)

func NextDate(now time.Time, date string, repeat string) (string, error) {
	ruleType := repeat[0:1]
	dateTime, err := time.Parse("20060102", date)
	if err != nil {
		return "", err
	}
	switch ruleType {
	case "d":
		value := repeat[2:]
		days, err := strconv.Atoi(value)
		if err != nil {
			return "", err
		}
		nextDate := dateTime.AddDate(0, 0, days)
		for now.After(nextDate) {
			nextDate = nextDate.AddDate(0, 0, days)
		}
		return nextDate.Format("20060102"), nil
	case "y":
		nextDate := dateTime.AddDate(1, 0, 0)
		for now.After(nextDate) {
			nextDate = nextDate.AddDate(1, 0, 0)
		}
		return nextDate.Format("20060102"), nil
	}
	return date, nil
}
