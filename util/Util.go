package util
import (
"time"
"strconv"
"fmt"
	"errors"
)

// returns the date formatted as mmddyyyy
func GetDate() string {
	t := time.Now()
	return t.Format("01022006")
}

//given a number as a string, increment it
func IncStringNumber(number string) (string, error) {
	i, err := strconv.Atoi(number)
	if err != nil {
		fmt.Printf("thisis num", string(i))
		return "", err
	} else {
		return strconv.Itoa(i+1), nil
	}
}

//given a number as a string, increment it
func CalcNewAverage(oldAverageStr string, numRoundsStr string, newScoreStr string) (string, error) {
	oldAverage, errA := strconv.ParseFloat(oldAverageStr, 64)
	newScore, errB := strconv.ParseFloat(newScoreStr, 64)
	numRounds, errC := strconv.ParseFloat(numRoundsStr, 64)

	if errA != nil || errB != nil || errC != nil {
		return "", errors.New("parameters were not numbers")
	}

	newAverage := (oldAverage * (numRounds - 1) + newScore) / numRounds

	return strconv.FormatFloat(newAverage, 'f', 10, 64), nil
}

func MakeTimestamp() string {
	return strconv.FormatInt(time.Now().UnixNano() / int64(time.Millisecond, ), 10)
}