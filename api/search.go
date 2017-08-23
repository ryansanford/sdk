package api

import (
	"net/http"
)

// Enum for Return Types.
type SearchType string

const (
	FileString        SearchType = "file"
	AcquisitionString SearchType = "acquisition"
	SessionString     SearchType = "session"
	AnalysisString    SearchType = "analysis"
)

// A single search query made to the API
type SearchQuery struct {
	ReturnType SearchType `json:"return_type"` // REQUIRED file|acquisition|session|analysis

	SearchString string        `json:"search_string,omitempty"` // OPTIONAL KEY any string including spaces and special characters
	AllData      bool          `json:"all_data,omitempty"`      // OPTIONAL, DEFAULTS TO FALSE true|false
	Filters      []interface{} `json:"filters,omitempty"`       // A LIST OF ES FILTERS, OPTIONAL KEY, find list of available filters here: https://www.elastic.co/guide/en/elasticsearch/reference/current/term-level-queries.html
	Size         string        `json:"size,omitempty"`          // OPTIONAL KEY if it is all, all files/other containers are returned
}

type ProjectSearchResponse struct {
	Id   string `json:"_id,omitempty"`
	Name string `json:"label,omitempty"`
}
type GroupSearchResponse struct {
	Id   string `json:"_id,omitempty"`
	Name string `json:"label,omitempty"`
}

// Runnning into Parsing errors when using time.Time
type SessionSearchResponse struct {
	Id        string `json:"_id,omitempty"`
	Archived  bool   `json:"archived,omitempty"`
	Name      string `json:"label,omitempty"`
	Timestamp string `json:"timestamp,omitempty"`
	Created   string `json:"created,omitempty"`
}
type AcquisitionSearchResponse struct {
	Id        string `json:"_id,omitempty"`
	Archived  bool   `json:"archived,omitempty"`
	Name      string `json:"label,omitempty"`
	Timestamp string `json:"timestamp,omitempty"`
	Created   string `json:"created,omitempty"`
}
type SubjectSearchResponse struct {
	Code string `json:"code,omitempty"`
}
type FileSearchResponse struct {
	Measurements []string `json:"measurements,omitempty"`
	Created      string   `json:"created,omitempty"`
	Type         string   `json:"type,omitempty"`
	Name         string   `json:"name,omitempty"`
	Size         int      `json:"size,omitempty"`
}
type AnalysisSearchResponse struct {
	Id      string `json:"_id,omitempty"`
	Name    string `json:"label,omitempty"`
	User    string `json:"user,omitempty"`
	Created string `json:created,omitempty"`
}
type ParentSearchResponse struct {
	Type string `json:"type,omitempty"`
	Id   string `json:"_id,omitempty"`
}

// SourceResponse for the SearchResponse
type SourceResponse struct {
	Project     *ProjectSearchResponse     `json:"project,omitempty"`
	Group       *GroupSearchResponse       `json:"group,omitempty"`
	Session     *SessionSearchResponse     `json:"session,omitempty"`
	Acquisition *AcquisitionSearchResponse `json:"acquisition,omitempty"`
	Subject     *SubjectSearchResponse     `json:"subject,omitempty"`
	File        *FileSearchResponse        `json:"file,omitempty"`
	Permissions []*Permission              `json:"permissions,omitempty"`
	Analysis    *AnalysisSearchResponse    `json:"analysis,omitempty"`
	Parent      *ParentSearchResponse      `json:"parent,omitempty"`
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


func (c *Client) SearchRaw(search_query *SearchQuery) (*SearchResponseList, *http.Response, error) {
	var aerr *Error
	var response *SearchResponseList

	resp, err := c.New().Post("dataexplorer/search").BodyJSON(search_query).Receive(&response, &aerr)

	return response, resp, Coalesce(err, aerr)
}
func (c *Client) Search(search_query *SearchQuery) ([]*SourceResponse, *http.Response, error) {
	var response *SearchResponseList
	var cleanList []*SourceResponse

	response, http, err := c.SearchRaw(search_query)
	for _, result := range response.Results {
		cleanList = append(cleanList, result.Source)
	}
	return cleanList, http, err
}
