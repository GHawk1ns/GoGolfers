package model

// A small identifier for a golfer, not meant to be large or contain skill data
type Golfer struct {
	GolferId	string	`json:"golferId"`
	Name    	string  `json:"name"`
	ImageUrl	string	`json:"imgUrl"`
}