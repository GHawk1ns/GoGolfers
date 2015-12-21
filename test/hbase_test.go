package test

import (
	"testing"
	"github.com/ghawk1ns/golf/util"
	"github.com/ghawk1ns/golf/database"
	"github.com/ghawk1ns/golf/model"
)

func TestHBase(t *testing.T) {
	config := util.GetConfig()
	database.InitHBase(config.HBaseConfig)
}

func TestGeneralPutGet(t *testing.T) {
	err := database.Test("helloworld")
	if err != nil {
		t.Error(err)
	}
}

func TestGolferTotalRound(t *testing.T) {

	err := database.SetGolferNumRounds("1", "23")

	if err != nil {
		t.Error(err)
	}

	getResult, err := database.GetGolferNumRounds("1")

	if err != nil {
		t.Error(err)
	} else if getResult != "23" {
		println("result " + getResult)
		t.Fail()
	}
}


func TestPutRound(t *testing.T) {

	var round model.Round
	round.Date = util.GetDate()

	var scoreA = model.Score{"foo","107"}
	var scoreB = model.Score{"guy","87"}
	var scoreC = model.Score{"goku","64"}
	var scoreD = model.Score{"kimJongUn","18"}
	round.Scores = []model.Score{scoreA, scoreB, scoreC, scoreD}

	err := database.PutRound(round)

	if err != nil {
		t.Error(err)
	}
}

func TestGetScoresForGolfer(t *testing.T) {
	val, err := database.GetScoresForGolfer("kimJongUn")
	if err != nil {
		t.Error(err)
	} else if val == nil {
		t.Fail()
	} else if val[util.GetDate()] != "18" {
		t.Fail()
	}
}


func TestGolferAverage(t *testing.T) {

	err := database.SetGolferAverage("average_golfer", util.GetDate())

	if err != nil {
		t.Error(err)
	}

	getResult, err := database.GetGolferAverage("average_golfer")

	if err != nil {
		t.Error(err)
	} else if getResult != util.GetDate() {
		t.Fail()
	}
}

func TestGetStats(t *testing.T) {

	database.SetGolferAverage("stats_man", "24")
	database.SetGolferNumRounds("stats_man", "56")

	roundAvg := make(chan string)
	numRounds := make(chan string)

	go func() {
		result, err := database.GetGolferAverage("stats_man")
		if err != nil {
			t.Error(err)
		} else {
			roundAvg <- result
		}
	}()

	go func() {
		result, err := database.GetGolferNumRounds("stats_man")
		if err != nil {
			t.Error(err)
		} else {
			numRounds <- result
		}
	}()

	stats := model.Stats{ <- numRounds, <- roundAvg, nil}

	if stats.Average != "24" {
		t.Fail()
	}

	if stats.Rounds != "56" {
		t.Fail()
	}
}