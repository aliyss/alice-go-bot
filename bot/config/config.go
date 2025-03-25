package config

import (
	"os"
)

type Config struct {
	Token   string
	AppID   string
	GuildID string
}

func GetConfig() Config {

	TOKEN := os.Getenv("TOKEN")
	if TOKEN == "" {
		panic("No token found")
	}

	APP_ID := os.Getenv("APP_ID")
	if APP_ID == "" {
		panic("No app id found")
	}

	GUILD_ID := os.Getenv("GUILD_ID")
	if GUILD_ID == "" {
		panic("No guild id found")
	}

	return Config{
		Token:   TOKEN,
		AppID:   APP_ID,
		GuildID: GUILD_ID,
	}
}
