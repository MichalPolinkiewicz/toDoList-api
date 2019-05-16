package main

import (
	"github.com/MichalPolinkiewicz/to-do-api/auth"
	"github.com/MichalPolinkiewicz/to-do-api/logger"
	"github.com/MichalPolinkiewicz/to-do-api/routes"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	router := newRouter()
	go log.Fatal(http.ListenAndServe(":8081", router))
}

func newRouter() *mux.Router {
	router := mux.NewRouter()

	for _, route := range routes.AllRoutes {
		var handler http.Handler

		handler = route.HandlerFunc

		//logger decorator
		handler = logger.Log(handler, route.Name)

		//auth decorator
		if route.Name != "SignIn" {
			handler = auth.CheckJwtToken(handler)
		}

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}
	return router
}
