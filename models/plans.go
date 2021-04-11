package models

import "gorm.io/gorm"

type Plan struct {
	gorm.Model
	Name     string
	Speed    string
	Duration int // Number of days
	Cost     int
	Notes    string //Additoinal string
}
