package api

import (
	"net/http"
)

// Enum for job states.
type SearchType string

const (
	FileString        SearchType = "file"
	AcquisitionString SearchType = "acquisition"
	SessionString     SearchType = "session"
	AnalysisString    SearchType = "analysis"
)

// A single search query made to the API
type SearchQuery struct {
	ReturnType string `json:"return_type"` // REQUIRED file|acquisition|session|analysis

	SearchString string   `json:"search_string,omitempty"` // OPTIONAL KEY any string including spaces and special characters
	AllData      bool     `json:"all_data,omitempty"`      // OPTIONAL, DEFAULTS TO FALSE true|false
	Filters      []string `json:"filters,omitempty"`       // A LIST OF ES FILTERS, OPTIONAL KEY, find list of available filters here: https://www.elastic.co/guide/en/elasticsearch/reference/current/term-level-queries.html
	Size         string   `json:"size,omitempty"`          // OPTIONAL KEY if it is all, all files/other containers are returned
}

// SourceResponse for the SearchResponse
type SourceResponse struct {
	Project     map[string]interface{} `json:"project,omitempty"`
	Group       map[string]interface{} `json:"group,omitempty"`
	Session     map[string]interface{} `json:"session,omitempty"`
	Acquisition map[string]interface{} `json:"acquisition,omitempty"`
	Subject     map[string]interface{} `json:"subject,omitempty"`
	File        map[string]interface{} `json:"file,omitempty"`
	Permissions []*Permission          `json:"permissions,omitempty"`
	Analysis    map[string]interface{} `json:"analysis,omitempty"`
}

// SearchResponse is used for endpoints of data_explorer
type SearchResponse struct {
	Id     string          `json:"_id"`
	Source *SourceResponse `json:"_source,omitempty"`
}

// Because the endpoint returns a key results which is a list of responses
type SearchResponseList struct {
	Results []*SearchResponse `json:"results,omitempty"`
}

func (c *Client) Search(search_query *SearchQuery) (*SearchResponseList, *http.Response, error) {
	var aerr *Error
	var response *SearchResponseList

	resp, err := c.New().Post("dataexplorer/search").BodyJSON(search_query).Receive(&response, &aerr)

	return response, resp, Coalesce(err, aerr)
}
