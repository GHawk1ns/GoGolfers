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
	"github.com/ghawk1ns/golf/blah"
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

	var round model.Round
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		handleError(w, err, 422)
	} else if err := r.Body.Close(); err != nil {
		handleError(w, err, 500)
	} else if err := json.Unmarshal(body, &round); err != nil {
		handleError(w, err, 400)
	} else if err := validateRound(round); err != nil {
		handleError(w, err, 400)
	} else {
		if err := database.PutRound(round); err != nil {
			handleError(w, err, 500)
		} else {
			go updateStats(round)
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.WriteHeader(http.StatusCreated)
			if err := json.NewEncoder(w).Encode(round); err != nil {
				panic(err)
			}
		}
	}
}

func handleError(w http.ResponseWriter, err error, code int) {
	blah.Error.Println("Error occured ", err.Error())
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(code) // unprocessable entity
	if err := json.NewEncoder(w).Encode(err.Error()); err != nil {
		panic(err)
	}
}

func validateRound(round model.Round) error {
	blah.Info.Println("Validating Round")
	if round.Date == "" {
		round.Date = util.GetDate()
	}

	if len(round.Scores) == 0 {
		return errors.New("Must include at least 1 score")
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
	blah.Info.Println("Round is Valid")
	return nil
}

/**
	Updates stats after a new round has been submitted
 */
func updateStats(round model.Round) {
	for _, score := range round.Scores {
		golferId := score.GolferId
		score := score.Score
		blah.Info.Printf("Updating stats for %s who just shot a %d\n", golferId, score)

		numRounds, err := database.GetGolferNumRounds(golferId)

		if err != nil {
			blah.Error.Println(err.Error())
			continue
		} else {
			numRounds++
			blah.Info.Printf("%s: new numRounds: %s\n", golferId, numRounds)
		}

		err = database.SetGolferNumRounds(golferId, numRounds)

		if err != nil {
			blah.Error.Println(err.Error())
			continue
		}

		currentAverage, err := database.GetGolferAverage(golferId)

		if err != nil {
			blah.Error.Println(err.Error())
			continue
		} else {
			blah.Info.Printf("%s: currentAverage: %s\n", golferId, currentAverage)
		}

		newAverage := util.CalcNewAverage(currentAverage, numRounds, score)
		blah.Info.Printf("%s: newAverage: %s\n", golferId, newAverage)

		err = database.SetGolferAverage(golferId, newAverage)
		if err != nil {
			blah.Error.Println(err.Error())
		}

		// Increase victory count over other golfers
		wins, err := database.GetGolferWins(golferId)
		if err != nil {
			blah.Error.Println(err.Error())
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