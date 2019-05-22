package handlers

import (
	"encoding/json"
	"github.com/MichalPolinkiewicz/to-do-api/auth"
	"github.com/MichalPolinkiewicz/to-do-api/db"
	"github.com/MichalPolinkiewicz/to-do-api/models"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func CreateTask(w http.ResponseWriter, req *http.Request) {
	var newTask models.Task
	_ = json.NewDecoder(req.Body).Decode(&newTask)
	newTask.UserId = auth.GetUserIdFromRequest(req)

	if newTask.IsValidTask() && newTask.UserId != 0 {
		db.CreateTask(&newTask)
		json.NewEncoder(w).Encode(db.GetAllTasks(&newTask.UserId))
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

func GetAllTasks(w http.ResponseWriter, req *http.Request) {
	id := auth.GetUserIdFromRequest(req)
	json.NewEncoder(w).Encode(db.GetAllTasks(&id))
}

func GetTaskById(w http.ResponseWriter, req *http.Request) {
	task := models.Task{}

	if id, ok := mux.Vars(req)["id"]; ok {
		idAsInt, _ := strconv.Atoi(id)
		task = db.GetTaskById(&idAsInt)
	}

	json.NewEncoder(w).Encode(task)
}

func GetTasksByStatus(w http.ResponseWriter, req *http.Request) {
	reqPrms := mux.Vars(req)
	var tasks []models.Task

	if status, ok := reqPrms["status"]; ok {
		statusAsInt, _ := strconv.Atoi(status)
		tasks = *getTasksByStatus(&statusAsInt)
	}
	json.NewEncoder(w).Encode(tasks)
}

func getTasksByStatus(s *int) *[]models.Task {
	return db.GetTasksByStatus(s)
}
