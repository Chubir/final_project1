package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

func CreateTaskHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		var data createRequest
		var respDataId createResponseId
		var respDataError createResponseError
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			respDataError.Error = err.Error()
			json.NewEncoder(w).Encode(&respDataError)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if data.Date == "" {
			respDataError.Error = "Не заполнено обязательное поле дата."
			json.NewEncoder(w).Encode(&respDataError)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if data.Title == "" {
			respDataError.Error = "Не заполнено обязательное поле задача."
			json.NewEncoder(w).Encode(&respDataError)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		// обротился за помощью к чату GPT
		re := regexp.MustCompile(`^(d\s\d+|y|w\s[1-7](,\s?[1-7])*)$`)
		if !re.MatchString(data.Repeat) && data.Repeat != "" {
			respDataError.Error = "Формат повторения задач не верный."
			json.NewEncoder(w).Encode(&respDataError)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		timeFutureTask, err := time.Parse("20060102", data.Date)
		if err != nil || timeFutureTask.Before(time.Now().Add(-24*time.Hour)) {
			respDataError.Error = "Введена дата несоответствующего формата."
			json.NewEncoder(w).Encode(&respDataError)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		query := "INSERT INTO scheduler (date, title, comment, repeat) VALUES ($1, $2, $3, $4)"
		res, err := db.Exec(query, data.Date, data.Title, data.Comment, data.Repeat)
		if err != nil {
			respDataError.Error = err.Error()
			json.NewEncoder(w).Encode(&respDataError)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		id, _ := res.LastInsertId()
		respDataId.Id = strconv.Itoa(int(id))
		json.NewEncoder(w).Encode(&respDataId)
		w.WriteHeader(http.StatusOK)
	}
}

func ListTaskHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		var respData listErrResponse
		tasks := []task{}
		query := "SELECT * FROM scheduler ORDER BY date"
		rows, err := db.Query(query)
		if err != nil {
			respData.Error = err.Error()
			json.NewEncoder(w).Encode(&respData)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		for rows.Next() {
			var rowTask task
			err := rows.Scan(&rowTask.Id, &rowTask.Date, &rowTask.Title, &rowTask.Comment, &rowTask.Repeat)
			if err != nil {
				respData.Error = err.Error()
				json.NewEncoder(w).Encode(&respData)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			tasks = append(tasks, rowTask)
		}
		var listResponse listResponse
		listResponse.Tasks = tasks
		json.NewEncoder(w).Encode(listResponse)
		w.WriteHeader(http.StatusOK)
	}
}

func ReadTaskHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		var respData listErrResponse
		id := r.URL.Query().Get("id")
		if id == "" {
			respData.Error = "Не указан идентификатор."
			json.NewEncoder(w).Encode(&respData)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		query := "SELECT * FROM scheduler WHERE id=$1"
		row := db.QueryRow(query, id)
		if row.Err() == sql.ErrNoRows {
			respData.Error = "Задача не найдена."
			json.NewEncoder(w).Encode(&respData)
			w.WriteHeader(http.StatusNotFound)
			return
		} else if row.Err() != nil {
			respData.Error = row.Err().Error()
			json.NewEncoder(w).Encode(&respData)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var rowTask task
		err := row.Scan(&rowTask.Id, &rowTask.Date, &rowTask.Title, &rowTask.Comment, &rowTask.Repeat)
		if err != nil {
			respData.Error = err.Error()
			json.NewEncoder(w).Encode(&respData)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		json.NewEncoder(w).Encode(rowTask)
		w.WriteHeader(http.StatusOK)

	}
}

func UpdateTaskHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		var data task
		var respData listErrResponse
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			respData.Error = err.Error()
			json.NewEncoder(w).Encode(&respData)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if data.Date == "" {
			respData.Error = "Не заполнено обязательное поле дата."
			json.NewEncoder(w).Encode(&respData)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if data.Title == "" {
			respData.Error = "Не заполнено обязательное поле задача."
			json.NewEncoder(w).Encode(&respData)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		timeFutureTask, err := time.Parse("20060102", data.Date)
		if err != nil || timeFutureTask.Before(time.Now()) {
			respData.Error = "Введена дата несоответствующего формата."
			json.NewEncoder(w).Encode(&respData)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		query := "UPDATE scheduler SET date = $1, title = $2, comment = $3, repeat = $4 WHERE id = $5"
		_, err = db.Exec(query, data.Date, data.Title, data.Comment, data.Repeat, data.Id)
		if err != nil {
			respData.Error = err.Error()
			json.NewEncoder(w).Encode(&respData)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
	}

}
func DoneTaskHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var respData listErrResponse
		id := r.URL.Query().Get("id")
		if id == "" {
			respData.Error = "Не указан идентификатор."
			json.NewEncoder(w).Encode(&respData)
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(http.StatusNotFound)
			return
		}
		query := "SELECT * FROM scheduler WHERE id=$1"
		row := db.QueryRow(query, id)
		if row.Err() == sql.ErrNoRows {
			respData.Error = "Задача не найдена"
			json.NewEncoder(w).Encode(&respData)
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(http.StatusNotFound)
			return
		} else if row.Err() != nil {
			respData.Error = row.Err().Error()
			json.NewEncoder(w).Encode(&respData)
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var rowTask task
		err := row.Scan(&rowTask.Id, &rowTask.Date, &rowTask.Title, &rowTask.Comment, &rowTask.Repeat)
		if err != nil {
			respData.Error = err.Error()
			json.NewEncoder(w).Encode(&respData)
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if rowTask.Repeat == "" {
			query := "DELETE FROM scheduler WHERE id = $1"
			_, err = db.Exec(query, rowTask.Id)
			if err != nil {
				respData.Error = err.Error()
				json.NewEncoder(w).Encode(&respData)
				w.Header().Set("Content-Type", "application/json; charset=UTF-8")
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		} else {
			newDate, err := NextDate(time.Now(), rowTask.Date, rowTask.Repeat)
			if err != nil {
				respData.Error = err.Error()
				json.NewEncoder(w).Encode(&respData)
				w.Header().Set("Content-Type", "application/json; charset=UTF-8")
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			query := "UPDATE scheduler SET date = $1 WHERE id = $2"
			_, err = db.Exec(query, newDate, rowTask.Id)
			if err != nil {
				respData.Error = err.Error()
				json.NewEncoder(w).Encode(&respData)
				w.Header().Set("Content-Type", "application/json; charset=UTF-8")
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}
		w.WriteHeader(http.StatusOK)
	}
}

func DeleteTaskHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		var respData listErrResponse
		id := r.URL.Query().Get("id")
		if id == "" {
			respData.Error = "Не указан идентификатор."
			json.NewEncoder(w).Encode(&respData)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		query := "DELETE FROM scheduler WHERE id = $1"
		res, err := db.Exec(query, id)
		if err != nil {
			respData.Error = err.Error()
			json.NewEncoder(w).Encode(&respData)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		count, _ := res.RowsAffected()
		if count == 0 {
			respData.Error = "Неверный идентификатор."
			json.NewEncoder(w).Encode(&respData)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
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
		nextTime, err := time.Parse("20060102", next)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Формат даты next неверный."))
			return
		}

		if nextTime.Before(nowTime) {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Введена дата несоответствующего формата."))
			return
		}
		nextDate, err := NextDate(nowTime, next, repeat)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Ошибка в функции NextDate."))
			return
		}
		w.Write([]byte(nextDate))
	}
}
