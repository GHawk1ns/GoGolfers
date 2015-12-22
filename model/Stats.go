package model

// Describes the skill level of a golfer
type Stats struct {
	Rounds    int     	 `json:"rounds"`
	Average   float64    `json:"average"`
	Wins 	  []WinCount `json:"wins"`
}