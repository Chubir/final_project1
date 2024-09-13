package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

func NextDate(now time.Time, repeat string, date string) (string, error) {
	now, _ = time.Parse("20060102", now.Format("20060102"))
	ruleType := repeat[0:1]
	dateTime, err := time.Parse("20060102", date)
	if err != nil {
		return "", err
	}
	switch ruleType {
	case "d":
		value := repeat[2:]
		days, err := strconv.Atoi(value)
		if days > 400 {
			return "", fmt.Errorf("нельзя больше 400")
		}
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

func SendBadRequest(w http.ResponseWriter, err error) {
	rErr := ResponseError{Error: err.Error()}
	err = json.NewEncoder(w).Encode(rErr)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusBadRequest)
}

func SendNotFound(w http.ResponseWriter, err error) {
	rErr := ResponseError{Error: err.Error()}
	err = json.NewEncoder(w).Encode(rErr)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNotFound)
}
