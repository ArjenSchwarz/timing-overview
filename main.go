package main

import (
	"context"
	"flag"
	"os"
	"time"

	"github.com/ArjenSchwarz/TimingSDK/config"
	"github.com/ArjenSchwarz/TimingSDK/parser"
	"github.com/aws/aws-lambda-go/lambda"
)

func HandleRequest(ctx context.Context) (string, error) {
	// configuration := config.ParseEnvironmentConfig()
	return "test", nil
}

var local bool
var rawstartdate string
var rawenddate string

func init() {
	flag.BoolVar(&local, "local", false, "Run locally")
	flag.StringVar(&rawstartdate, "startdate", "", "The startdate in format 2006-01-02 15:04")
	flag.StringVar(&rawenddate, "enddate", "", "The enddate, leave blank for now")
	flag.Parse()
}

func main() {
	if local {
		configuration := config.ParseConfigFile()
		if rawstartdate != "" {
			configuration.StartDate, _ = time.ParseInLocation("2006-01-02 15:04", rawstartdate, time.Local)
		}
		if rawenddate != "" {
			configuration.EndDate, _ = time.ParseInLocation("2006-01-02 15:04", rawenddate, time.Local)
		}
		f, _ := os.Create("output.png")
		defer f.Close()
		parser.CreateProjectOverviewPieChart(configuration, f)
	} else {
		lambda.Start(HandleRequest)
	}

}
