package config

import (
	"context"
	"encoding/json"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
)

// Configuration is a configuration object
type Configuration struct {
	// The start date for any calls
	StartDate time.Time
	// The end date for any calls
	EndDate time.Time
	// The API token for Timing
	Token string `json:"APIToken"`
}

// ParseConfigFile will parse a local config.json file for its configuration
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

// ParseEnvironmentConfig retrieves configuration from the environment
// This is used for the Lambda only and will expect the API_TOKEN_PARAMETER
// environment variable to be set
func ParseEnvironmentConfig() Configuration {
	ctx := context.TODO()
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		panic(err)
	}
	svc := ssm.NewFromConfig(cfg)
	parameter, err := svc.GetParameter(ctx, &ssm.GetParameterInput{
		Name:           aws.String(os.Getenv("API_TOKEN_PARAMETER")),
		WithDecryption: aws.Bool(true),
	})
	if err != nil {
		panic(err)
	}
	configuration := Configuration{
		Token:     *parameter.Parameter.Value,
		StartDate: time.Now(),
		EndDate:   time.Now(),
	}
	return configuration
}

// ParseJSON parses a JSON document provided
// The JSON format is expected to have the start_date and end_date keys.
func (configuration *Configuration) ParseJSON(jsonstring string) {
	type requestBody struct {
		StartDate string `json:"start_date"`
		EndDate   string `json:"end_date"`
	}
	body := requestBody{}
	err := json.Unmarshal([]byte(jsonstring), &body)
	if err != nil {
		panic(err)
	}
	startdate, err := time.Parse("2006-01-02 15:04 -0700", body.StartDate)
	if err != nil {
		panic(err)
	}
	configuration.StartDate = startdate
	enddate, err := time.Parse("2006-01-02 15:04 -0700", body.EndDate)
	if err != nil {
		panic(err)
	}
	configuration.EndDate = enddate
}
