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

		if route.Name != "SignIn" && route.Name != "Logout" {
			handler = logger.Log(handler, route.Name)
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
