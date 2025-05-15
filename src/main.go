package main

import (
	"bufio"
	"context"
	"encoding/base64"
	"flag"
	"io"
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
	content, _ := io.ReadAll(reader)
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

func runLocal() {
	configuration := config.ParseConfigFile()
	if rawstartdate != "" {
		configuration.StartDate, _ = time.ParseInLocation("2006-01-02 15:04", rawstartdate, time.Local)
	} else {
		now := time.Now()
		year, month, day := now.Date()
		configuration.StartDate = time.Date(year, month, day, 0, 0, 0, 0, time.Local)
	}
	if rawenddate != "" {
		configuration.EndDate, _ = time.ParseInLocation("2006-01-02 15:04", rawenddate, time.Local)
	}
	f, _ := os.Create("output.png")
	defer f.Close()
	parser.CreateProjectOverviewPieChart(configuration, f)
}

var local bool
var rawstartdate string
var rawenddate string

func init() {
	_, taskExists := os.LookupEnv("LAMBDA_TASK_ROOT")
	_, execExists := os.LookupEnv("AWS_EXECUTION_ENV")
	local = true
	if taskExists || execExists {
		local = false
	}
	flag.StringVar(&rawstartdate, "startdate", "", "The startdate in format 2006-01-02 15:04, leave empty for start of today")
	flag.StringVar(&rawenddate, "enddate", "", "The enddate, leave empty for current date and time")
	flag.Parse()
}

func main() {
	if local {
		runLocal()
	} else {
		lambda.Start(handleRequest)
	}
}
