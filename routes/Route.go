package routes

import (
	"net/http"
	"github.com/gorilla/mux"
	"github.com/ghawk1ns/golf/blah"
)

func NewRouter() *mux.Router {

	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler

		handler = route.HandlerFunc
		handler = blah.HTTPLogger(handler, route.Name)

		router.
		Methods(route.Method).
		Path(route.Pattern).
		Name(route.Name).
		Handler(handler)
	}

	return router
}


