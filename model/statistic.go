package model

import "gorm.io/gorm"

type Statistic struct {
	gorm.Model
	AlgorithmName string
	BaseCurrency  string
	TrackCurrency string
	Exchange      string
	Result        string
	Action        string
}
