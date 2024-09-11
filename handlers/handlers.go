package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"myApp/internal/repo"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

func CreateTaskHandler(db *repo.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		var data createRequest
		var respDataId createResponseId
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			SendBadRequest(w, err)
			return
		}

		if data.Date != "" {
			date, err := time.Parse("20060102", data.Date)
			if err != nil {
				err = fmt.Errorf("неверный формат поля дата")
				SendBadRequest(w, err)
				return
			}
			now, _ := time.Parse("20060102", time.Now().Format("20060102"))
			if now.After(date) && data.Repeat == "" {
				data.Date = now.Format("20060102")
			} else if now.After(date) {
				newDate, err := NextDate(now, data.Repeat, data.Date)
				if err != nil {
					err = fmt.Errorf("не сработала функция NextDate")
					SendBadRequest(w, err)
					return
				}
				data.Date = newDate
			}
		} else {
			data.Date = time.Now().Format("20060102")
		}

		if data.Title == "" {
			err := fmt.Errorf("незаполнено обязательное поле задача")
			SendBadRequest(w, err)
			return
		}
		// обротился за помощью к чату GPT
		re := regexp.MustCompile(`^(d\s\d+|y|w\s[1-7](,\s?[1-7])*)$`)
		if !re.MatchString(data.Repeat) && data.Repeat != "" {
			err := fmt.Errorf("формат повторения задач не верный")
			SendBadRequest(w, err)
			return
		}

		timeFutureTask, err := time.Parse("20060102", data.Date)
		if err != nil || timeFutureTask.Before(time.Now().Add(-24*time.Hour)) {
			err := fmt.Errorf("введена дата несоответствующего формата")
			SendBadRequest(w, err)
			return
		}

		id, err := db.InsertTask(data.Date, data.Title, data.Comment, data.Repeat)
		if err != nil {
			SendBadRequest(w, err)
			return
		}
		respDataId.Id = strconv.Itoa(id)
		json.NewEncoder(w).Encode(respDataId)
		w.WriteHeader(http.StatusOK)
	}
}

func ListTaskHandler(db *repo.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		tasks, err := db.GetListTask()
		if err != nil {
			SendBadRequest(w, err)
			return
		}
		var listResponse listResponse
		listResponse.Tasks = tasks
		json.NewEncoder(w).Encode(listResponse)
		w.WriteHeader(http.StatusOK)
	}
}

func ReadTaskHandler(db *repo.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		id := r.URL.Query().Get("id")
		if id == "" {
			err := fmt.Errorf("не указан идентификатор")
			SendNotFound(w, err)
			return
		}

		task, err := db.GetTask(id)
		if err == sql.ErrNoRows {
			err := fmt.Errorf("задача не найдена")
			SendNotFound(w, err)
			return
		} else if err != nil {
			SendBadRequest(w, err)
			return
		}
		json.NewEncoder(w).Encode(task)
		w.WriteHeader(http.StatusOK)

	}
}

func UpdateTaskHandler(db *repo.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		var data task
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			SendBadRequest(w, err)
			return
		}

		re := regexp.MustCompile(`^(d\s\d+|y|w\s[1-7](,\s?[1-7])*)$`)
		if !re.MatchString(data.Repeat) && data.Repeat != "" {
			err := fmt.Errorf("формат повторения задач не верный")
			SendBadRequest(w, err)
			return
		}

		if data.Date == "" {
			err := fmt.Errorf("не заполнено обязательное поле дата")
			SendBadRequest(w, err)
			return
		}
		if data.Title == "" {
			err := fmt.Errorf("не заполнено обязательное поле задача")
			SendBadRequest(w, err)
			return
		}

		_, err = time.Parse("20060102", data.Date)
		if err != nil {
			err := fmt.Errorf("введена дата несоответствующего формата")
			SendBadRequest(w, err)
			return
		}

		task, err := db.UpdateTask(data.Date, data.Title, data.Comment, data.Repeat, data.Id)
		if err != nil {
			SendBadRequest(w, err)
			return
		}

		json.NewEncoder(w).Encode(task)
		w.WriteHeader(http.StatusOK)
	}

}
func DoneTaskHandler(db *repo.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		if id == "" {
			err := fmt.Errorf("не указан идентификатор")
			SendNotFound(w, err)
			return
		}

		task, err := db.GetTask(id)
		if err == sql.ErrNoRows {
			err := fmt.Errorf("задача не найдена")
			SendNotFound(w, err)
			return
		} else if err != nil {
			SendBadRequest(w, err)
			return
		}

		if task.Repeat == "" {
			err := db.DeleteTask(id)
			if err != nil {
				SendBadRequest(w, err)
				return
			}
		} else {
			newDate, err := NextDate(time.Now(), task.Repeat, task.Date)
			if err != nil {
				SendBadRequest(w, err)
				return
			}
			err = db.UpdateDate(newDate, task.Id)
			if err != nil {
				SendBadRequest(w, err)
				return
			}
		}
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		json.NewEncoder(w).Encode(nil)
		w.WriteHeader(http.StatusOK)
	}
}

func DeleteTaskHandler(db *repo.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		id := r.URL.Query().Get("id")
		if id == "" {
			err := fmt.Errorf("не указан идентификатор")
			SendNotFound(w, err)
			return
		}
		err := db.DeleteTask(id)
		if err != nil {
			SendBadRequest(w, err)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		json.NewEncoder(w).Encode(nil)
		w.WriteHeader(http.StatusOK)

	}
}

func NextDateHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		now := r.URL.Query().Get("now")
		next := r.URL.Query().Get("date")
		repeat := r.URL.Query().Get("repeat")
		re := regexp.MustCompile(`^(d\s\d+|y|w\s[1-7](,\s?[1-7])*)$`)
		if repeat == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Отсутсвует повторение задачи."))
			return
		}
		if !re.MatchString(repeat) {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Формат повторения задач неверный."))
			return
		}

		nowTime, err := time.Parse("20060102", now)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Формат даты now неверный."))
			return
		}

		if next == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("date обязательное поле."))
		}

		_, err = time.Parse("20060102", next)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Формат даты next неверный."))
			return
		}

		nextDate, err := NextDate(nowTime, repeat, next)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Ошибка в функции NextDate."))
			return
		}
		w.Write([]byte(nextDate))
	}
}
