package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"
)

func CreateTaskHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		var data createRequest
		var respData createResponse
		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			respData.Error = err.Error()
			json.NewEncoder(w).Encode(&respData)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		query := "INSERT INTO scheduler (date, title, comment, repeat) VALUES ($1, $2, $3, $4)"
		res, err := db.Exec(query, data.Date, data.Title, data.Comment, data.Repeat)
		if err != nil {
			respData.Error = err.Error()
			json.NewEncoder(w).Encode(&respData)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		id, _ := res.LastInsertId()
		respData.Id = int(id)
		json.NewEncoder(w).Encode(&respData)
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
			respData.Error = "Не указан идентификатор"
			json.NewEncoder(w).Encode(&respData)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		query := "SELECT * FROM scheduler WHERE id=$1"
		row := db.QueryRow(query, id)
		if row.Err() == sql.ErrNoRows {
			respData.Error = "Задача не найдена"
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
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		var respData listErrResponse
		id := r.URL.Query().Get("id")
		if id == "" {
			respData.Error = "Не указан идентификатор"
			json.NewEncoder(w).Encode(&respData)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		query := "SELECT * FROM scheduler WHERE id=$1"
		row := db.QueryRow(query, id)
		if row.Err() == sql.ErrNoRows {
			respData.Error = "Задача не найдена"
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
		if rowTask.Repeat == "" {
			query := "DELETE FROM scheduler WHERE id = $1"
			_, err = db.Exec(query, rowTask.Id)
			if err != nil {
				respData.Error = err.Error()
				json.NewEncoder(w).Encode(&respData)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		} else {
			newDate, err := NextDate(time.Now(), rowTask.Date, rowTask.Repeat)
			if err != nil {
				respData.Error = err.Error()
				json.NewEncoder(w).Encode(&respData)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			query := "UPDATE scheduler SET date = $1 WHERE id = $2"
			_, err = db.Exec(query, newDate, rowTask.Id)
			if err != nil {
				respData.Error = err.Error()
				json.NewEncoder(w).Encode(&respData)
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
			respData.Error = "Не указан идентификатор"
			json.NewEncoder(w).Encode(&respData)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		query := "DELETE FROM scheduler WHERE id = $1"
		_, err := db.Exec(query, id)
		if err != nil {
			respData.Error = err.Error()
			json.NewEncoder(w).Encode(&respData)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)

	}
}
