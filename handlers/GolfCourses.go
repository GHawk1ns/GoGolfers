package handlers

import (
	"fmt"
	"net/http"
	"encoding/json"
	"github.com/ghawk1ns/golf/database"
	"github.com/ghawk1ns/golf/logger"
)

func GolfCourses(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	rawCourses, err := database.GetCourses()
	if err != nil {
		onGolferError(w, err)
	} else {
		b, err := json.Marshal(rawCourses)
		if err != nil {
			onCourseError(w, err)
		} else {
			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, string(b))
		}
	}
}

func onCourseError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	logger.Error.Println("an error occured:", err)
	fmt.Fprintln(w, nil)
}