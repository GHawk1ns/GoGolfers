package routes

import (
	"net/http"
	"github.com/ghawk1ns/golf/handlers"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/",
		handlers.Index,
	},
	Route{
		"Index",
		"GET",
		"/test",
		handlers.Index,
	},
	Route{
		"RoundCreate",
		"POST",
		"/roundCreate",
		handlers.RoundCreate,
	},
	Route{
		"Golfers",
		"GET",
		"/golfers",
		handlers.Golfers,
	},
	Route{
		"GolferProfile",
		"GET",
		"/golfers/{id}",
		handlers.GolferProfile,
	},
	Route{
		"GolfCourses",
		"GET",
		"/courses",
		handlers.GolfCourses,
	},
}