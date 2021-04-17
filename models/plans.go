package models

type Plan struct {
	ID       int    `json:"id,string"`
	Name     string `json:"name"`
	Speed    string `json:"speed"`
	Duration int    `json:"duration"` // Number of days
	Cost     int    `json:"cost"`
	Notes    string `json:"notes"` //Additoinal string
	IsActive bool   `json:"isactive"`
}
