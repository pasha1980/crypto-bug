package config

import "github.com/joho/godotenv"

func environmentSetup()  {
	_ = godotenv.Load()
}
