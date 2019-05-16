package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/MichalPolinkiewicz/to-do-api/db"
	"github.com/MichalPolinkiewicz/to-do-api/models"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func CreateTask(w http.ResponseWriter, req *http.Request){
	var newTask models.Task
	_ = json.NewDecoder(req.Body).Decode(&newTask)

	if newTask.IsValidTask(){
		db.CreateTask(&newTask)
	} else {
		fmt.Println("Task is invalid!")
	}
	json.NewEncoder(w).Encode(db.GetAllTasks())
}

func GetAllTasks(w http.ResponseWriter, req *http.Request) {
	json.NewEncoder(w).Encode(db.GetAllTasks())
}

func GetTaskById(w http.ResponseWriter, req *http.Request) {
	task := models.Task{}

	if id, ok := mux.Vars(req)["id"]; ok {
		idAsInt, _ := strconv.Atoi(id)
		task = db.GetTaskById(&idAsInt)
	}

	json.NewEncoder(w).Encode(task)
}

func GetTasksByStatus(w http.ResponseWriter, req *http.Request){
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
