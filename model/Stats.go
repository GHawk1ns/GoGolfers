package model

// Describes the skill level of a golfer
type Stats struct {
	Rounds    string     `json:"rounds"`
	Average   string     `json:"average"`
	Wins 	  []WinCount `json:"wins"`
}