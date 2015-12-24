package model

// Describes a round of golf
type Round struct {
	Date    	string  `json:"date"`
	CourseId    string  `json:"courseId"`
	Scores		[]Score `json:"scores"`
}