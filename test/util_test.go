package test

import (
"github.com/ghawk1ns/golf/util"
"testing"
"fmt"
)


func TestIncNumber(t *testing.T) {
	one := "456"
	two, err := util.IncStringNumber(one)
	fmt.Println(two)
	if err != nil {
		t.Error(err)
	} else if two != "457" {
		t.Fail()
	}
}

func TestCalcNewAvg(t *testing.T) {
	val, err := util.CalcNewAverage("2", "4", "4")
	if err != nil {
		t.Error(err)
	} else if val != "2.5000000000" {

		println(val)
		t.Fail()
	}
}
