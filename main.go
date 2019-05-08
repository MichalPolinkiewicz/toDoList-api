package main

import (
	"github.com/MichalPolinkiewicz/to-do-api/models"
	"github.com/MichalPolinkiewicz/to-do-api/routes"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	models.Tasks = []models.Task{
		{Id: "1", Name:"Task 1", Description:"Desc for tasc 1", Status:1},
		{Id: "2", Name:"Task 2", Description:"Desc for tasc 2", Status:2},
		{Id: "3", Name:"Task 3", Description:"Desc for tasc 3", Status:3},
	}

	router := newRouter()
	log.Fatal(http.ListenAndServe(":8081", router))

}

func newRouter() *mux.Router {
	router := mux.NewRouter()

	for _, route := range routes.AllRoutes {
		var handler http.Handler

		handler = route.HandlerFunc
		//handler = Logger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	return router
}
