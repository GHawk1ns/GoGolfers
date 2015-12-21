package handlers

import (
	"fmt"
	"net/http"
)

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome to golf showdown - check out /golfers and /golfers/<id>")
}