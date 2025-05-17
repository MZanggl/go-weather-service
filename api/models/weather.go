package models

import "gorm.io/gorm"

type Weather struct {
	gorm.Model
	RecordedAt  string
	Humidity    float64
	Temperature float64
}

func (w Weather) TableName() string {
	return "weather"
}
