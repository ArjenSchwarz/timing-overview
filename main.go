package main

import (
	"bufio"
	"context"
	"encoding/base64"
	"flag"
	"io/ioutil"
	"os"
	"time"

	"github.com/ArjenSchwarz/timing-overview/config"
	"github.com/ArjenSchwarz/timing-overview/parser"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	configuration := config.ParseEnvironmentConfig()
	body, err := base64.StdEncoding.DecodeString(request.Body)
	if err != nil {
		panic(err)
	}
	configuration.ParseJSON(string(body))
	f, _ := os.Create("/tmp/output.png")
	parser.CreateProjectOverviewPieChart(configuration, f)
	f.Close()
	f, _ = os.Open("/tmp/output.png")
	reader := bufio.NewReader(f)
	content, _ := ioutil.ReadAll(reader)
	// Encode as base64.
	encoded := base64.StdEncoding.EncodeToString(content)
	headers := make(map[string]string)
	headers["content-type"] = "image/png"
	return events.APIGatewayProxyResponse{
		Body:            encoded,
		Headers:         headers,
		StatusCode:      200,
		IsBase64Encoded: true,
	}, nil
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
		lambda.Start(handleRequest)
	}
}
