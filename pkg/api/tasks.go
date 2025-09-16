package api

import (
	"Go-final/pkg/db"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type TaskResponse struct {
	Tasks []db.Task `json:"tasks"`
}

func getTasksHandler(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")
	tasks, err := db.Tasks(search, 50)
	if err != nil {
		writeJson(w, http.StatusInternalServerError, map[string]string{"error": "task retrieval error"})
		return
	}
	writeJson(w, http.StatusOK, TaskResponse{Tasks: tasks})
}

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJson(w, http.StatusBadRequest, map[string]string{"error": "invalid id"})
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeJson(w, http.StatusNotFound, map[string]string{"error": "task not found"})
		return
	}
	log.Printf("task received")
	writeJson(w, http.StatusOK, task)
}
func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		writeJson(w, http.StatusBadRequest, map[string]string{"error": "invalid json body"})
		return
	}
	if task.Title == "" {
		writeJson(w, http.StatusInternalServerError, map[string]string{"error": "task title is required"})
		return
	}
	err = checkDate(&task)
	if err != nil {
		writeJson(w, http.StatusInternalServerError, map[string]string{"error": "date check error"})
		return
	}
	err = db.UpdateTask(&task)
	if err != nil {
		writeJson(w, http.StatusInternalServerError, map[string]string{"error": "task update error"})
		return
	}
	writeJson(w, http.StatusOK, map[string]string{})
}
func doneTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		log.Printf("invalid ID")
		writeJson(w, http.StatusBadRequest, map[string]string{"error": "id is required"})
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		log.Printf("faild getting task: %v", err)
		writeJson(w, http.StatusNotFound, map[string]string{"error": "task not found"})
		return
	}

	if task.Repeat == "" {
		if err := db.DeleteTask(id); err != nil {
			writeJson(w, http.StatusOK, struct{}{})
			return
		}
		writeJson(w, http.StatusOK, struct{}{})

	} else {
		now := time.Now()
		nextDate, err := NextDate(now, task.Date, task.Repeat)

		if err != nil {
			writeJson(w, http.StatusBadRequest, map[string]string{"error": "invalid repeat"})
			return
		}

		if err = db.UpdateDate(nextDate, id); err != nil {
			writeJson(w, http.StatusOK, struct{}{})
			return
		}
		writeJson(w, http.StatusOK, struct{}{})
	}
}

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJson(w, http.StatusBadRequest, map[string]string{"error": "id is required"})
		return
	}
	if err := db.DeleteTask(id); err != nil {
		writeJson(w, http.StatusInternalServerError, map[string]string{"error": "task delete error"})
		return
	}
	writeJson(w, http.StatusOK, map[string]string{})
}
