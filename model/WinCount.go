package model

// Describes the number of wins a golfer has over an opponent
type WinCount struct {
	Golfer    Golfer    `json:"opponent"`
	Wins   	  int       `json:"wins"`
}

