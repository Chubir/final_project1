package repo

import (
	"database/sql"
	"fmt"
	"myApp/config"
	"os"

	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/sqlite3"
	"github.com/golang-migrate/migrate/source/file"
)

type Repository struct {
	db *sql.DB
}

func New(conf *config.Config) (*Repository, error) {
	db, err := sql.Open(config.SQL_TYPE, conf.Dbfile)
	if err != nil {
		panic(err)
	}
	repo := Repository{db: db}
	var install bool
	if _, err := os.Stat(conf.Dbfile); err != nil {
		install = true
	}

	if install {
		instance, err := sqlite3.WithInstance(db, &sqlite3.Config{})
		if err != nil {
			return nil, err
		}

		fSrc, err := (&file.File{}).Open("./migrations")
		if err != nil {
			return nil, err
		}

		m, err := migrate.NewWithInstance("file", fSrc, "scheduler", instance)
		if err != nil {
			return nil, err
		}

		// modify for Down
		if err := m.Up(); err != nil {
			return nil, err
		}
	}
	return &repo, nil
}

func (r *Repository) InsertTask(date, title, comment, repeat string) (int, error) {
	query := "INSERT INTO scheduler (date, title, comment, repeat) VALUES ($1, $2, $3, $4)"
	res, err := r.db.Exec(query, date, title, comment, repeat)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	return int(id), err
}

func (r *Repository) GetListTask() ([]Task, error) {
	tasks := []Task{}
	query := "SELECT * FROM scheduler ORDER BY date"
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var rowTask Task
		err := rows.Scan(&rowTask.Id, &rowTask.Date, &rowTask.Title, &rowTask.Comment, &rowTask.Repeat)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, rowTask)
	}
	if rows.Err() != nil {
		return nil, err
	}
	return tasks, err
}

func (r *Repository) GetTask(id string) (Task, error) {

	query := "SELECT * FROM scheduler WHERE id=$1"
	row := r.db.QueryRow(query, id)

	if row.Err() != nil {
		return Task{}, row.Err()
	}

	var rowTask Task
	err := row.Scan(&rowTask.Id, &rowTask.Date, &rowTask.Title, &rowTask.Comment, &rowTask.Repeat)

	return rowTask, err

}

func (r *Repository) UpdateTask(date, title, comment, repeat, id string) (Task, error) {

	query := "UPDATE scheduler SET date = $1, title = $2, comment = $3, repeat = $4 WHERE id = $5"
	result, err := r.db.Exec(query, date, title, comment, repeat, id)
	if err != nil {
		return Task{}, err
	}

	if count, _ := result.RowsAffected(); count == 0 {
		err = fmt.Errorf("задача не была обновлена")
		return Task{}, err
	}
	task := Task{Date: date, Title: title, Comment: comment, Repeat: repeat, Id: id}
	return task, err
}

func (r *Repository) DeleteTask(id string) error {
	query := "DELETE FROM scheduler WHERE id = $1"
	res, err := r.db.Exec(query, id)
	count, _ := res.RowsAffected()
	if count == 0 {
		err = fmt.Errorf("неверный идентификатор")
		return err
	}
	return err
}

func (r *Repository) UpdateDate(date, id string) error {
	query := "UPDATE scheduler SET date = $1 WHERE id = $2"
	_, err := r.db.Exec(query, date, id)
	return err
}
