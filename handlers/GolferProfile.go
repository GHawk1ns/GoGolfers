package handlers

import (
	"fmt"
	"net/http"
	"encoding/json"
	"github.com/ghawk1ns/golf/database"
	"github.com/gorilla/mux"
	"github.com/ghawk1ns/golf/model"
	"errors"
	"github.com/ghawk1ns/golf/logger"
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

	scores, err := database.GetScoresForGolfer(golfer.Id)
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
		logger.Info.Println("golferProfile: ", string(b))
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, string(b))
	}
}

// This is not a good way to do this
func onGolferProfileError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusBadRequest)
	logger.Error.Println("an error occured:", err.Error())
	fmt.Fprintln(w, nil)
}

func getStats(golferId string) (model.Stats, error) {

	roundAvg, err := database.GetAllAveragesForGolfer(golferId)
	if err != nil {
		logger.Error.Println(err)
		roundAvg = nil
	}

	// retrieve the golfer's total rounds played
	numRounds, err := database.GetAllRoundsForGolfer(golferId)
	if err != nil {
		logger.Error.Println(err)
		numRounds = nil
	}

	// retrieve the golfer's victory over other golfers
	logger.Info.Println("Getting win stats for", golferId)
	wins, err := database.GetGolferWins(golferId)
	var winCounts []model.WinCount
	if err != nil {
		logger.Error.Println(err)
		winCounts = nil
	} else {
		var localWinCounts []model.WinCount
		for opponentId,count := range wins {
			golfer, err := database.GetGolferById(opponentId)
			if err != nil {
				logger.Error.Println(err.Error())
			} else {
				logger.Info.Printf("%s has beaten %s, %d times\n", golferId, golfer.Name, count)
				localWinCounts = append(localWinCounts, model.WinCount{golfer, count})
			}

		}
		winCounts = localWinCounts
	}

	stats := model.Stats{ numRounds, roundAvg, winCounts}
	if stats.Rounds == nil || stats.Averages == nil {
		return stats, errors.New("something went wrong with stat gathering")
	} else {
		return stats, nil
	}
}