package api

import (
	"Go-final/pkg/db"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

func writeJson(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("writeJson error: %v", err)
	}
}

func checkDate(task *db.Task) error {
	now := time.Now()
	var next string
	if task.Date == "" {
		task.Date = now.Format(DateLayout)
	}
	t, err := time.Parse(DateLayout, task.Date)
	if err != nil {
		return fmt.Errorf("date format error: %w", err)
	}
	if len(task.Repeat) > 0 {
		next, err = NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return fmt.Errorf("next date error: %w", err)
		}
	}
	if t.After(now) {
		if task.Repeat == "" {
			task.Date = now.Format(DateLayout)
		} else {
			task.Date = next
		}
	} else {
		task.Date = now.Format(DateLayout)
	}
	return nil
}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task *db.Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		log.Printf("writeJson error: %v", err)
		writeJson(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON"})
		return
	}
	if task.Title == "" {
		log.Printf("title error: empty title")
		writeJson(w, http.StatusBadRequest, map[string]string{"error": "empty title"})
		return
	}
	err = checkDate(task)
	if err != nil {
		log.Printf("checkDate error: %v", err)
		writeJson(w, http.StatusBadRequest, map[string]string{"error": "invalid checkDate"})
		return
	}
	var id int64
	id, err = db.AddTask(task)
	if err != nil {
		log.Printf("AddTask error: %v", err)
		writeJson(w, http.StatusInternalServerError, map[string]string{"error": "failed to add task"})
		return
	}
	idStr := strconv.FormatInt(id, 10)
	writeJson(w, http.StatusCreated, map[string]string{"id": idStr})
}
