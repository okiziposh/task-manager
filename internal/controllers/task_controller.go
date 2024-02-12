package controllers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/project/task-manager/internal/models"
)

var TaskStore = make(map[int]*models.Task)

type TaskController struct {
	DB *sql.DB
}

func (tc *TaskController) CreateTask(w http.ResponseWriter, r *http.Request) {
	var task models.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		log.Println("Failed to decode request body:", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate task fields
	if task.Title == "" || task.Status == "" {
		http.Error(w, "Title and Status are required fields", http.StatusBadRequest)
		return
	}

	// Insert task into database
	result, err := tc.DB.Exec("INSERT INTO tasks (title, description, status) VALUES (?, ?, ?)", task.Title, task.Description, task.Status)
	if err != nil {
		log.Println("Failed to insert task into database:", err)
		http.Error(w, "Failed to create task", http.StatusInternalServerError)
		return
	}

	// Get ID of newly inserted task
	id, err := result.LastInsertId()
	if err != nil {
		log.Println("Failed to get ID of inserted task:", err)
		http.Error(w, "Failed to create task", http.StatusInternalServerError)
		return
	}

	task.ID = int(id)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

func (tc *TaskController) GetTasks(w http.ResponseWriter, r *http.Request) {
	tasks := []models.Task{}

	rows, err := tc.DB.Query("SELECT * FROM tasks")
	if err != nil {
		log.Println("Failed to retrieve tasks:", err)
		http.Error(w, "Failed to retrieve tasks", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var task models.Task
		err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.Status)
		if err != nil {
			log.Println("Failed to scan task row:", err)
			http.Error(w, "Failed to retrieve tasks", http.StatusInternalServerError)
			return
		}
		tasks = append(tasks, task)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func (tc *TaskController) GetTaskByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskID, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Println("Invalid task ID:", err)
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	var task models.Task
	err = tc.DB.QueryRow("SELECT * FROM tasks WHERE id = ?", taskID).Scan(&task.ID, &task.Title, &task.Description, &task.Status)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Task not found", http.StatusNotFound)
			return
		}
		log.Println("Failed to retrieve task:", err)
		http.Error(w, "Failed to retrieve task", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

func (tc *TaskController) UpdateTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskID, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Println("Invalid task ID:", err)
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	var task models.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		log.Println("Failed to decode request body:", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate task fields
	if task.Title == "" || task.Status == "" {
		http.Error(w, "Title and Status are required fields", http.StatusBadRequest)
		return
	}

	_, err = tc.DB.Exec("UPDATE tasks SET title = ?, description = ?, status = ? WHERE id = ?", task.Title, task.Description, task.Status, taskID)
	if err != nil {
		log.Println("Failed to update task:", err)
		http.Error(w, "Failed to update task", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (tc *TaskController) DeleteTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskID, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Println("Invalid task ID:", err)
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	_, err = tc.DB.Exec("DELETE FROM tasks WHERE id = ?", taskID)
	if err != nil {
		log.Println("Failed to delete task:", err)
		http.Error(w, "Failed to delete task", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
