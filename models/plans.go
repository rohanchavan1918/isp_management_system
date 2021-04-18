package models

import "time"

type Plan struct {
	ID       int    `json:"id,string"`
	Name     string `json:"name"`
	Speed    string `json:"speed"`
	Duration int    `json:"duration"` // Number of days
	Cost     int    `json:"cost"`
	Notes    string `json:"notes"` //Additoinal string
	IsActive bool   `json:"isactive"`
}

type UserPlans struct {
	ID         int  `json:"id,string"`
	UserId     int  `json:"user_id,string"`
	PlanId     int  `json:"plan_id,string"`
	IsActive   bool `json:"isactive"`
	ValidTill  time.Time
	Created_at time.Time
	UpdatedAt  time.Time
	DeletedAt  time.Time
}
