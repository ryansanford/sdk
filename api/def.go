package api

import (
	"time"
)

// Error is an API error. All failed server responses should be of this form.
// TODO: implement error interface, change coalesce
type Error struct {
	Message    string `json:"message"`
	StatusCode int    `json:"status_code"`
}

// Permission represents the capability of a single user on a given container. Many containers have an array of these permissions, and they are frequently casscaded down the container hierarchy.
type Permission struct {
	Id    string `json:"_id"`
	Level string `json:"access"`
}

type Origin struct {
	Id     string `json:"id,omitempty"`
	Method string `json:"method,omitempty"`
	Name   string `json:"name,omitempty"`
	Type   string `json:"type,omitempty"`
}

type ContainerReference struct {
	Id   string `json:"id,omitempty"`
	Type string `json:"type,omitempty"`
}

type FileReference struct {
	Id   string `json:"id,omitempty"`
	Type string `json:"type,omitempty"`
	Name string `json:"name,omitempty"`
}

// IdResponse is used for endpoints that respond with an Id.
type IdResponse struct {
	Id string `json:"_id"`
}

// ModifiedResponse is used for endpoints that respond with a count of modified objects.
type ModifiedResponse struct {
	ModifiedCount int `json:"modified"`
}

// SourceResponse for the SearchResponse
type SourceResponse struct {
	Project 	*Project 	 	`json:"project,omitempty"`
	Group		*Group		 	`json:"group,omitempty"`
	Session		*Session	 	`json:"session,omitempty"`
	Acquisition	*Acquisition 	`json:"acquisition,omitempty"`
	Subject		*Subject		`json:"subject,omitempty"`
	File 		*File 			`json:"file,omitempty"`
}

// SearchResponse is used for endpoints of data_explorer
type SearchResponse struct {
	Id 		string				`json:"_id"`
	Source	*SourceResponse		`json:"_source,omitempty"`
}

// DeleteResponse is used for endpoints that respond with a count of deleted objects.
type DeletedResponse struct {
	DeletedCount int `json:"deleted"`
}

type Note struct {
	Id     string `json:"id,omitempty"`
	UserId string `json:"user,omitempty"`
	Text   string `json:"text,omitempty"`

	Created  *time.Time `json:"created,omitempty"`
	Modified *time.Time `json:"modified,omitempty"`
}
