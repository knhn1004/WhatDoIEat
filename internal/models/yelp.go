package models

type YelpResponse struct {
	Businesses []Restaurant `json:"businesses"`
}
