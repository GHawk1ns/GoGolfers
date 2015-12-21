package handlers

import (
	"fmt"
	"net/http"
	"encoding/json"
	"github.com/ghawk1ns/golf/database"
	"github.com/ghawk1ns/golf/blah"
)

func Golfers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	rawGolfers, err := database.GetGolfers()
	if err != nil {
		onGolferError(w, err)
	} else {
		b, err := json.Marshal(rawGolfers)
		if err != nil {
			onGolferError(w, err)
		} else {
			w.WriteHeader(http.StatusOK)
			fmt.Fprintln(w, string(b))
		}
	}
}

func onGolferError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	blah.Error.Println("an error occured:", err)
	fmt.Fprintln(w, nil)
}