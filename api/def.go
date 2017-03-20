package api

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
	Site  string `json:"site,omitempty"`
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

// DeleteResponse is used for endpoints that respond with a count of deleted objects.
type DeletedResponse struct {
	DeletedCount int `json:"deleted"`
}
