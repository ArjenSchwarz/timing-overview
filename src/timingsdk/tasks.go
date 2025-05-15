package timingsdk

import (
	"encoding/json"
	"io"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

const entriesURL = "https://web.timingapp.com/api/v1/time-entries"

// TaskData is the top level task response object
type TaskData struct {
	Data    []TaskDetails `json:"data"`
	Message string        `json:"message"`
	Links   Links         `json:"links"`
	Meta    Meta          `json:"meta"`
}

// TaskDetails contains the details about the tasks
type TaskDetails struct {
	Self      string         `json:"self"`
	StartDate time.Time      `json:"start_date"`
	EndDate   time.Time      `json:"end_date"`
	Duration  float64        `json:"duration"`
	Project   ProjectDetails `json:"project"`
	Title     string         `json:"title"`
	Notes     string         `json:"notes"`
	IsRunning bool           `json:"is_running"`
}

// Links contains links for paging purposes
type Links struct {
	First    string `json:"first"`
	Last     string `json:"last"`
	Next     string `json:"next"`
	Previous string `json:"prev"`
}

// Meta contains the metadata for the request
type Meta struct {
	CurrentPage int64  `json:"current_page"`
	From        int64  `json:"from"`
	LastPage    int64  `json:"last_page"`
	Path        string `json:"path"`
	PerPage     int64  `json:"per_page"`
	To          int64  `json:"to"`
	Total       int64  `json:"total"`
}

// ProjectDetails contains the project level details for a task
type ProjectDetails struct {
	Self              string   `json:"self"`
	Title             string   `json:"title"`
	TitleChain        []string `json:"title_chain"`
	Color             string   `json:"color"`
	ProductivityScore int64    `json:"productivity_score"`
	IsArchived        bool     `json:"is_archived"`
	Parent            struct {
		Self string `json:"self"`
	} `json:"parent"`
}

// TaskGrouping is a grouped set of tasks
// For example, tasks can be grouped by task or project to make it easier to
// have a complete overview
type TaskGrouping struct {
	Title    string
	Duration float64
	Color    string
	Tasks    []TaskDetails
}

// TaskCollection contains a collection of grouped tasks
type TaskCollection struct {
	Title     string
	Duration  float64
	Groupings map[string]TaskGrouping
	keys      []string
}

// EntriesRequest contains the flags for the entries request
type EntriesRequest struct {
	// StartDate is the earliest time an entry can have started
	StartDate time.Time
	// EndDate is the latest time an entry can have started
	EndDate time.Time
	// Restricts the query to tasks associated with the given projects. Need to provided as project paths like /projects/12345
	Projects []string
	// Restricts the query to tasks whose title and/or notes contain all words in this parameter. The search is case-insensitive but diacritic-sensitive.
	SearchQuery string
	// If provided (as a string of "true" or "false"), returns only tasks that are either running or not running
	IsRunning string
	// If true, the properties of the task's project will be included in the response
	IncludeProjectData bool
	// If true, the response will also contain tasks that belong to any child projects of the ones provided in Projects.
	IncludeChildProjects bool
	// The authentication token
	Token string
}

// GetTimeEntries retrieves a TaskData object with task entries
func GetTimeEntries(request EntriesRequest) (TaskData, error) {
	result := TaskData{}
	client := &http.Client{}
	req, err := http.NewRequest("GET", entriesURL, nil)
	if err != nil {
		return result, err
	}
	q := req.URL.Query()
	if request.IsRunning != "" {
		q.Add("is_running", request.IsRunning)
	}
	if !request.StartDate.IsZero() {
		q.Add("start_date_min", request.StartDate.Format("2006-01-02T15:04:05-07:00"))
	}
	if !request.EndDate.IsZero() {
		q.Add("start_date_max", request.EndDate.Format("2006-01-02T15:04:05-07:00"))
	}
	if request.IncludeProjectData {
		q.Add("include_project_data", "true")
	}
	if request.IncludeChildProjects {
		q.Add("include_child_projects", "true")
	}
	if len(request.Projects) > 0 {
		for _, project := range request.Projects {
			q.Add("projects[]", project)
		}
	}
	if request.SearchQuery != "" {
		q.Add("search_query", request.SearchQuery)
	}
	req.URL.RawQuery = q.Encode()
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	req.Header.Add("Authorization", "Bearer "+request.Token)
	resp, err := client.Do(req)
	if err != nil {
		return result, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return result, err
	}
	json.Unmarshal(body, &result)
	return result, nil
}

// GroupEntriesByTask Groups a provided entries list by task
func GroupEntriesByTask(data TaskData) TaskCollection {
	result := TaskCollection{Title: "Grouped by task", Groupings: make(map[string]TaskGrouping)}
	for _, entry := range data.Data {
		result.AddEntry(entry, entry.Title)
	}
	return result
}

// GroupEntriesByProject Groups a provided entries list by Project
func GroupEntriesByProject(data TaskData) TaskCollection {
	result := TaskCollection{Title: "Grouped by project", Groupings: make(map[string]TaskGrouping)}
	for _, entry := range data.Data {
		result.AddEntry(entry, strings.Join(entry.Project.TitleChain, "/"))
	}
	return result
}

// AddEntry adds a TaskDetails object to the TaskGrouping
func (grouping *TaskGrouping) AddEntry(entry TaskDetails) {
	grouping.Tasks = append(grouping.Tasks, entry)
	grouping.Duration = grouping.Duration + entry.Duration
}

// GetDuration returns the summed duration of all tasks in the TaskGrouping
func (grouping *TaskGrouping) GetDuration() (time.Duration, error) {
	duration, err := time.ParseDuration(strconv.FormatFloat(grouping.Duration, 'f', 0, 64) + "s")
	if err != nil {
		return 0, err
	}
	return duration, nil
}

// AddEntry adds a TaskDetails object to a TaskCollection, while also placing it in the defined grouping
func (collection *TaskCollection) AddEntry(entry TaskDetails, groupTitle string) {
	grouping := TaskGrouping{Title: groupTitle, Color: entry.Project.Color}
	groupings := collection.Groupings
	if _, ok := groupings[groupTitle]; ok {
		grouping = groupings[groupTitle]
	} else {
		collection.keys = append(collection.keys, groupTitle)
		sort.Strings(collection.keys)
	}
	grouping.AddEntry(entry)
	collection.Groupings[groupTitle] = grouping
	collection.Duration = collection.Duration + entry.Duration
}

// GetDuration returns the summed duration of all tasks in the collection
func (collection *TaskCollection) GetDuration() (time.Duration, error) {
	duration, err := time.ParseDuration(strconv.FormatFloat(collection.Duration, 'f', 0, 64) + "s")
	if err != nil {
		return 0, err
	}
	return duration, nil
}

// GetGroupings returns the Tasks from the TaskCollection separated into predefined groupings
func (collection *TaskCollection) GetGroupings() []TaskGrouping {
	result := []TaskGrouping{}
	for _, k := range collection.keys {
		result = append(result, collection.Groupings[k])
	}
	return result
}
