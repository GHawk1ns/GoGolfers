package handlers

import (
	"errors"
	"encoding/json"
	"net/http"
	"io"
	"io/ioutil"
	"github.com/ghawk1ns/golf/model"
	"github.com/ghawk1ns/golf/database"
	"github.com/ghawk1ns/golf/util"
	"github.com/ghawk1ns/golf/logger"
	"fmt"
)

func RoundCreate(w http.ResponseWriter, r *http.Request) {

	secret := r.Header.Get("secret")
	if secret != util.GetSecret() {
		w.Header().Set("Content-Type", "application/text; charset=UTF-8")
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintln(w, "https://www.youtube.com/watch?v=QDySGUFAom0")
		return
	}

	round := &model.Round{}
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		handleError(w, err, 422)
	} else if err := r.Body.Close(); err != nil {
		handleError(w, err, 500)
	} else if err := json.Unmarshal(body, round); err != nil {
		handleError(w, err, 400)
	} else if err := validateRound(round); err != nil {
		handleError(w, err, 400)
	} else {
		if err := database.PutRound(*round); err != nil {
			handleError(w, err, 500)
		} else {
			go updateStats(*round)
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(http.StatusCreated)
			if err := json.NewEncoder(w).Encode(*round); err != nil {
				panic(err)
			}
		}
	}
}

func handleError(w http.ResponseWriter, err error, code int) {
	logger.Error.Println("Error occured ", err.Error())
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(code) // unprocessable entity
	if err := json.NewEncoder(w).Encode(err.Error()); err != nil {
		panic(err)
	}
}

func validateRound(round *model.Round) error {
	logger.Info.Println("Validating Round")
	if round.Date == "" {
		round.Date = util.GetDate()
	}

	if len(round.Scores) == 0 {
		  errors.New("Must include at least 1 score")
	} else if round.CourseId == "" {
		return errors.New("Must include course id")
	} else {
		for _, score := range round.Scores {
			if score.GolferId == "" {
				return errors.New("golferId must not be empty")
			}
			if score.Score == 0 {
				return errors.New("golferId: " + score.GolferId + " score must not be empty")
			}
		}
	}

	_, err := database.GetCourseById(round.CourseId)
	if err != nil {
		return errors.New("Invalid golf course id")
	}

	logger.Info.Println("Round is Valid")
	return nil
}

/**
	Updates stats after a new round has been submitted
 */
func updateStats(round model.Round) {

	// TODO: Update stats for golf course

	for _, score := range round.Scores {
		golferId := score.GolferId
		score := score.Score
		logger.Info.Printf("Updating stats for %s who just shot a %d\n", golferId, score)

		database.IncGolferTotalRounds(golferId)

		numRounds, err := database.GetGolferNumRounds(golferId, round.CourseId)

		if err != nil {
			logger.Error.Println(err.Error())
			continue
		} else {
			numRounds++
			logger.Info.Printf("%s: new numRounds: %s\n", golferId, numRounds)
		}

		err = database.SetGolferNumRounds(golferId, round.CourseId, numRounds)

		if err != nil {
			logger.Error.Println(err.Error())
			continue
		}

		currentAverage, err := database.GetGolferAverage(golferId, round.CourseId)

		if err != nil {
			logger.Error.Println(err.Error())
			continue
		} else {
			logger.Info.Printf("%s: currentAverage: %s\n", golferId, currentAverage)
		}

		newAverage := util.CalcNewAverage(currentAverage, numRounds, score)
		logger.Info.Printf("%s: newAverage: %s\n", golferId, newAverage)

		err = database.SetGolferAverage(golferId, round.CourseId, newAverage)
		if err != nil {
			logger.Error.Println(err.Error())
		}

		// Increase victory count over other golfers
		wins, err := database.GetGolferWins(golferId)
		if err != nil {
			logger.Error.Println(err.Error())
			continue
		}
		for _, opponentScoreInfo := range round.Scores {
			opponentId := opponentScoreInfo.GolferId
			if golferId != opponentId && score < opponentScoreInfo.Score {
				wins[opponentId] = wins[opponentId] + 1
			}
		}
		database.SetGolferWins(golferId, wins)
	}
}