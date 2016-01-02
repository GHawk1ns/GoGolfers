package model

// A small identifier for a golfer, not meant to be large or contain skill data
type Golfer struct {
	Id       string	`json:"id"`
	Name     string  `json:"name"`
	ImageUrl string	`json:"imgUrl"`
}