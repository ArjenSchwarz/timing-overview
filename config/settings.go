package config

import (
	"encoding/json"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
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
	mySession := session.Must(session.NewSession())
	// Create a SSM client from just a session.
	svc := ssm.New(mySession)
	parameter, err := svc.GetParameter(&ssm.GetParameterInput{
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

func (configuration *Configuration) ParseJson(jsonstring string) {
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
