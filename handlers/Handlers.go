package handlers

import (
	"encoding/json"
	"github.com/MichalPolinkiewicz/to-do-api/models"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func CreateTask(w http.ResponseWriter, req *http.Request){
	var newTask models.Task
	_ = json.NewDecoder(req.Body).Decode(&newTask)
	models.Tasks = append(models.Tasks, newTask)
	json.NewEncoder(w).Encode(models.Tasks)
}

func GetAllTasks(w http.ResponseWriter, req *http.Request) {
	json.NewEncoder(w).Encode(models.Tasks)
}

func GetTaskById(w http.ResponseWriter, req *http.Request) {
	reqPrms := mux.Vars(req)
	for _, task := range models.Tasks {
		if task.Id == reqPrms["id"] {
			json.NewEncoder(w).Encode(task)
			return
		}
	}
	json.NewEncoder(w).Encode(models.Task{})
}

func GetTasksByStatus(w http.ResponseWriter, req *http.Request){
	reqPrms := mux.Vars(req)
	var tasks []models.Task
	if status, ok := reqPrms["status"]; ok {
		statusAsInt, _ := strconv.ParseInt(status, 10, 64)
		tasks = getTasksByStatus(int(statusAsInt))
	}
	json.NewEncoder(w).Encode(tasks)
}

func getTasksByStatus(status int) []models.Task {
	var tasks []models.Task
	for _, task := range models.Tasks {
		if task.Status == status {
			tasks = append(tasks, task)
		}
	}
	return tasks
}
