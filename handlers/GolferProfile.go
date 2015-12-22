package handlers

import (
	"fmt"
	"net/http"
	"encoding/json"
	"github.com/ghawk1ns/golf/database"
	"github.com/gorilla/mux"
	"github.com/ghawk1ns/golf/model"
	"errors"
"github.com/ghawk1ns/golf/blah"
)

func GolferProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	vars := mux.Vars(r)
	golferId := vars["id"]

	if golferId == "" {
		onGolferError(w, errors.New("Invalid Golfer Id"))
	}

	var result model.Profile
	golfer, err := database.GetGolferById(golferId)
	if err != nil {
		onGolferError(w, err)
		return
	} else {
		result.Golfer = golfer
	}

	scores, err := database.GetScoresForGolfer(golfer.GolferId)
	if err != nil {
		onGolferError(w, err)
		return
	} else {
		result.Scores = scores
	}

	result.Stats, err = getStats(golferId)

	b, err := json.Marshal(result)
	if err != nil {
		onGolferProfileError(w, err)
		return
	} else {
		blah.Info.Println("golferProfile: ", string(b))
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, string(b))
	}
}

// This is not a good way to do this
func onGolferProfileError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusBadRequest)
	blah.Error.Println("an error occured:", err.Error())
	fmt.Fprintln(w, nil)
}

func getStats(golferId string) (model.Stats, error) {
	roundAvg := make(chan int)
	numRounds := make(chan int)

	go func() {
		result, err := database.GetGolferAverage(golferId)
		if err != nil {
			// TODO figure out wtf to do if something fails in a channel
			roundAvg <- -1
		} else {
			roundAvg <- result
		}
	}()

	go func() {
		result, err := database.GetGolferNumRounds(golferId)
		if err != nil {
			// TODO figure out wtf to do if something fails in a channel
			numRounds <- -1
		} else {
			numRounds <- result
		}
	}()

	stats := model.Stats{ <- numRounds, <- roundAvg, nil}
	if stats.Rounds == -1 || stats.Average == -1 {
		return stats, errors.New("something went wrong with stat gathering")
	} else {
		return stats, nil
	}
}