package test

import (
	"github.com/ghawk1ns/golf/util"
	"testing"
)

func TestCalcNewAvg(t *testing.T) {
	if util.CalcNewAverage(2, 4, 4) != 2.5000000000 {
		t.Fail()
	}
}
