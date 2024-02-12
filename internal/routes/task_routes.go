package routes

import (
	"database/sql"

	"github.com/gorilla/mux"
	"github.com/project/task-manager/internal/controllers"
)

func SetupRoutes(db *sql.DB) *mux.Router {
	router := mux.NewRouter()
	taskController := &controllers.TaskController{DB: db}

	// Create a new task
	router.HandleFunc("/tasks", taskController.CreateTask).Methods("POST")

	// Retrieve a list of tasks
	router.HandleFunc("/tasks", taskController.GetTasks).Methods("GET")

	// Retrieve a specific task by its ID
	router.HandleFunc("/tasks/{id}", taskController.GetTaskByID).Methods("GET")

	// Update a task
	router.HandleFunc("/tasks/{id}", taskController.UpdateTask).Methods("PUT")

	// Delete a task
	router.HandleFunc("/tasks/{id}", taskController.DeleteTask).Methods("DELETE")

	return router
}
