package models

import "gorm.io/gorm"

type Weather struct {
	gorm.Model
	RecordedAt  string
	Humidity    float32
	Temperature float32
}

func (w Weather) TableName() string {
	return "weather"
}
