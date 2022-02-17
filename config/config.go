package config

import "time"

var RestartSeconds time.Duration = 60

func Initialization() {
	environmentSetup()
	databaseSetup()
}
