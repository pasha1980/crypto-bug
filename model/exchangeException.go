package model

import "gorm.io/gorm"

type ExchangeException struct {
	gorm.Model
	Exchange      string
	BaseCurrency  string
	TrackCurrency string
}
