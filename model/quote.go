package model

import (
	"gorm.io/gorm"
	"time"
)

type Quote struct {
	gorm.Model
	Exchange         string
	Date             time.Time
	BaseCurrency     string
	TrackingCurrency string
	Value            float64
}
