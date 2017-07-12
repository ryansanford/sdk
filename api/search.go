package api

import (
	"net/http"
)

// A single search query made to the API
type SearchQuery struct {
	ReturnType string `json:"return_type,omitempty"` // REQUIRED file|acquisition|session|analysis

	SearchString string   `json:"search_string,omitempty"` // OPTIONAL KEY any string including spaces and special characters
	AllData      *bool    `json:"all_data,omitempty"`      // OPTIONAL, DEFAULTS TO FALSE true|false
	Filters      []string `json:"filters,omitempty"`       // A LIST OF ES FILTERS, OPTIONAL KEY, find list of available filters here: https://www.elastic.co/guide/en/elasticsearch/reference/current/term-level-queries.html
}

func (c *Client) Search(search_query *SearchQuery) (*SearchResponseList, *http.Response, error) {
	var aerr *Error
	var response *SearchResponseList

	resp, err := c.New().Post("dataexplorer/search").BodyJSON(search_query).Receive(&response, &aerr)

	return response, resp, Coalesce(err, aerr)
}
