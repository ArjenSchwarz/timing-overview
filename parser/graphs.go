package parser

import (
	"fmt"
	"io"

	"github.com/ArjenSchwarz/timing-overview/config"
	"github.com/ArjenSchwarz/timing-overview/timingsdk"
	"github.com/wcharczuk/go-chart"
	"github.com/wcharczuk/go-chart/drawing"
)

// CreateProjectOverviewPieChart generates a pie chart of all the projects run in the provided timespan
func CreateProjectOverviewPieChart(configuration config.Configuration, target io.Writer) {
	request := timingsdk.EntriesRequest{
		Token:              configuration.Token,
		StartDate:          configuration.StartDate,
		EndDate:            configuration.EndDate,
		IncludeProjectData: true,
	}
	test, err := timingsdk.GetTimeEntries(request)
	if err != nil {
		panic(err)
	}
	chartValues := []chart.Value{}

	tasks := timingsdk.GroupEntriesByProject(test)
	if err != nil {
		panic(err)
	}
	for _, taskgroup := range tasks.GetGroupings() {
		duration, err := taskgroup.GetDuration()
		if err != nil {
			panic(err)
		}
		// fmt.Printf("%v: %v\n", taskgroup.Title, duration.String())
		// for _, task := range taskgroup.Tasks {
		// 	fmt.Printf("%v: %v - %v - %v: %v\n", task.Title, task.Project.Color, task.StartDate.Local().Format("Mon, 02 Jan 2006 15:04:05 MST"), task.EndDate.Local().Format("Mon, 02 Jan 2006 15:04:05 MST"), task.Duration)
		// }
		if taskgroup.Color != "" {
			value := chart.Value{
				Value: taskgroup.Duration,
				Label: fmt.Sprintf("%v (%v)", taskgroup.Title, duration.String()),
				Style: chart.Style{
					FillColor: drawing.ColorFromHex(taskgroup.Color[1:7]),
					FontSize:  12,
				},
			}
			chartValues = append(chartValues, value)
		}
	}
	pie := chart.PieChart{
		Title:  fmt.Sprintf("Data for %v", request.StartDate.Format("02 Jan 2006")),
		Width:  768,
		Height: 768,
		Values: chartValues,
	}
	pie.Render(chart.PNG, target)
}
