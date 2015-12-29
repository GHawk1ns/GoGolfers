package model

// Describes the skill level of a golfer
type Stats struct {
	Rounds    map[string]int		`json:"rounds"`
	Averages  map[string]float64 	`json:"averages"`
	Wins 	  []WinCount			`json:"wins"`
}