package util
import (
	"time"
	"strconv"
)


//given a number as a string, increment it
func CalcNewAverage(oldAverage float64, numRounds int, newScore int) float64 {
	return (oldAverage * (float64(numRounds) - 1) + float64(newScore)) / float64(numRounds)
}

func IntToStr(num int) string {
	return strconv.Itoa(num)
}

func FloatToStr(num float64) string {
	return strconv.FormatFloat(num, 'f', 10, 64)
}

func StrToInt(num string) (int, error) {
	return strconv.Atoi(num)
}

func StrToFloat(num string) (float64, error) {
	return strconv.ParseFloat(num, 64)
}

func MakeTimestamp() string {
	return strconv.FormatInt(time.Now().UnixNano() / int64(time.Millisecond, ), 10)
}
// returns the date formatted as mmddyyyy
func GetDate() string {
	t := time.Now()
	return t.Format("01022006")
}
