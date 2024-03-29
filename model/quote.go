package model

import (
	"gorm.io/gorm"
	"time"
)

type Quote struct {
	gorm.Model
	Exchange      string
	Date          time.Time
	BaseCurrency  string
	TrackCurrency string
	Value         float64
	IsAbnormally  bool `gorm:"default:0;not null"`
}
