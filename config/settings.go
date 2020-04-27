package config

import (
	"encoding/json"
	"os"
	"time"
)

// Config is a configuration object
type Configuration struct {
	// The start date for any calls
	StartDate time.Time
	// The end date for any calls
	EndDate time.Time
	// The API token for Timing
	Token string `json:"APIToken"`
}

func ParseConfigFile() Configuration {
	file, _ := os.Open("config.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err := decoder.Decode(&configuration)
	if err != nil {
		panic(err)
	}
	if configuration.StartDate.IsZero() {
		configuration.StartDate = time.Now().Local()
	}
	if configuration.EndDate.IsZero() {
		configuration.EndDate = time.Now().Local()
	}
	return configuration
}

func ParseEnvironmentConfig() Configuration {
	configuration := Configuration{
		Token:     os.Getenv("API_TOKEN"),
		StartDate: time.Now(),
		EndDate:   time.Now(),
	}
	return configuration
}
