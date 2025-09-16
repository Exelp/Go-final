package db

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func AddTask(task *Task) (int64, error) {

	q := `INSERT INTO scheduler (date,title,comment,repeat) VALUES (?,?,?,?)`
	res, err := db.Exec(q, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, fmt.Errorf("add task: %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("add task: %w", err)
	}
	return id, err
}

func Tasks(search string, limit int) ([]Task, error) {
	var rows *sql.Rows
	var err error

	t, err := time.Parse("02.01.2006", search)
	if err == nil {
		q := `SELECT id,date,title,comment,repeat FROM scheduler WHERE date = ? LIMIT ?`
		rows, err = db.Query(q, t.Format("20060102"), limit)
	} else if search != "" {
		likeSearch := "%" + search + "%"
		q := `SELECT id,date,title,comment,repeat FROM scheduler WHERE title LIKE ? OR comment LIKE ? ORDER BY date LIMIT ?`
		rows, err = db.Query(q, likeSearch, likeSearch, limit)
	} else {
		q := `SELECT id,date,title,comment,repeat FROM scheduler ORDER BY date LIMIT ?`
		rows, err = db.Query(q, limit)
	}
	if err != nil {
		return nil, fmt.Errorf("error get tasks: %w", err)
	}
	defer rows.Close()
	var tasks []Task

	for rows.Next() {
		var task Task
		if err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
			return nil, fmt.Errorf("error scan tasks: %w", err)
		}
		tasks = append(tasks, task)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error get tasks: %w", err)
	}
	if tasks == nil {
		tasks = []Task{}
	}
	return tasks, nil
}

func GetTask(id string) (*Task, error) {
	var task Task
	q := `SELECT id,date,title,comment,repeat FROM scheduler WHERE id = ?`
	err := db.QueryRow(q, id).Scan(
		&task.ID,
		&task.Date,
		&task.Title,
		&task.Comment,
		&task.Repeat)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("task not found")
	}
	if err != nil {
		return nil, fmt.Errorf("error get task: %w", err)
	}
	return &task, nil
}

func UpdateTask(task *Task) error {
	q := `UPDATE scheduler SET date=?,title=?,comment=?,repeat=? WHERE id=?`
	res, err := db.Exec(q, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return fmt.Errorf("error update task: %w", err)
	}
	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("error RowsAffected update: %w", err)
	}
	if count == 0 {
		return errors.New("task not found")
	}
	return nil
}
func UpdateDate(next string, id string) error {
	query := `UPDATE scheduler SET date = :date WHERE id = :id`
	res, err := db.Exec(query,
		sql.Named("date", next),
		sql.Named("id", id),
	)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf(`incorrect id for updating task`)
	}
	return nil
}

func DeleteTask(id string) error {
	q := `DELETE FROM scheduler WHERE id=?`
	res, err := db.Exec(q, id)
	if err != nil {
		return fmt.Errorf("error delete: %w", err)
	}
	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("error RowsAffected delete: %w", err)
	}
	if count == 0 {
		return errors.New("task not found")
	}
	return nil
}
