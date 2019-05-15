package main

import (
	"github.com/MichalPolinkiewicz/to-do-api/logger"
	"github.com/MichalPolinkiewicz/to-do-api/routes"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {

	//models.Tasks = db.GetAllTasks()
	router := newRouter()
	go log.Fatal(http.ListenAndServe(":8081", router))

}

func newRouter() *mux.Router {
	router := mux.NewRouter()

	for _, route := range routes.AllRoutes {
		var handler http.Handler

		handler = route.HandlerFunc

		//decorator
		handler = logger.Log(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}
	return router
}
