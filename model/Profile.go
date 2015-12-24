package model

// The complete package describing a golfer
type Profile struct {
	Golfer  Golfer						`json:"golfer"`
	Stats 	Stats 	 				  	`json:"stats"`
	Scores  map[string]map[string]int   `json:"scores"`
}