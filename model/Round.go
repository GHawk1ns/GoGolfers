package model

// Describes a round of golf
type Round struct {
	Date    string  `json:"date"`
	Scores	[]Score `json:"scores"`
}