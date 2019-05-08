package routes

import (
	"github.com/MichalPolinkiewicz/to-do-api/handlers"
	"net/http"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var AllRoutes = Routes{
	Route{
		"CreateTask",
		"POST",
		"/task",
		handlers.CreateTask,
	},
	Route{
		"GetAllTasks",
		"GET",
		"/tasks",
		handlers.GetAllTasks,
	},
	Route{
		"GetTaskById",
		"GET",
		"/task/{id}",
		handlers.GetTaskById,
	},
	Route{
		"GetTasksByStatus",
		"GET",
		"/tasks/{status}",
		handlers.GetTasksByStatus,
	},
}