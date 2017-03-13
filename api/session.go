package api

import (
	"errors"
	"net/http"
	"strconv"
	"time"
)

type Subject struct {
	Id   string `json:"_id,omitempty"`
	Name string `json:"code,omitempty"`
}

type Session struct {
	Id        string `json:"_id,omitempty"`
	Name      string `json:"label,omitempty"`
	GroupId   string `json:"group,omitempty"`
	ProjectId string `json:"project,omitempty"`

	Subject   *Subject   `json:"subject,omitempty"`
	Timestamp *time.Time `json:"timestamp,omitempty"`

	Created  *time.Time `json:"created,omitempty"`
	Modified *time.Time `json:"modified,omitempty"`
	Files    []*File    `json:"files,omitempty"`

	Public      bool          `json:"public,omitempty"`
	Permissions []*Permission `json:"permissions,omitempty"`
}

func (c *Client) GetAllSessions() ([]*Session, *http.Response, error) {
	var aerr *Error
	var sessions []*Session
	resp, err := c.New().Get("sessions").Receive(&sessions, &aerr)
	return sessions, resp, Coalesce(err, aerr)
}

func (c *Client) GetSession(id string) (*Session, *http.Response, error) {
	var aerr *Error
	var session *Session
	resp, err := c.New().Get("sessions/"+id).Receive(&session, &aerr)
	return session, resp, Coalesce(err, aerr)
}

func (c *Client) AddSession(session *Session) (string, *http.Response, error) {
	var aerr *Error
	var response *IdResponse
	var result string

	resp, err := c.New().Post("sessions").BodyJSON(session).Receive(&response, &aerr)

	if response != nil {
		result = response.Id
	}

	return result, resp, Coalesce(err, aerr)
}

func (c *Client) ModifySession(id string, session *Session) (*http.Response, error) {
	var aerr *Error
	var response *ModifiedResponse

	resp, err := c.New().Put("sessions/"+id).BodyJSON(session).Receive(&response, &aerr)

	// Should not have to check this count
	// https://github.com/scitran/core/issues/680
	if err == nil && aerr == nil && response.ModifiedCount != 1 {
		return resp, errors.New("Modifying session " + id + " returned " + strconv.Itoa(response.ModifiedCount) + " instead of 1")
	}

	return resp, Coalesce(err, aerr)
}

func (c *Client) DeleteSession(id string) (*http.Response, error) {
	var aerr *Error
	var response *DeletedResponse

	resp, err := c.New().Delete("sessions/"+id).Receive(&response, &aerr)

	// Should not have to check this count
	// https://github.com/scitran/core/issues/680
	if err == nil && aerr == nil && response.DeletedCount != 1 {
		return resp, errors.New("Deleting session " + id + " returned " + strconv.Itoa(response.DeletedCount) + " instead of 1")
	}

	return resp, Coalesce(err, aerr)
}

func (c *Client) UploadToSession(id string, files ...*UploadSource) (chan int64, chan error) {

	url := "sessions/"+id+"/files"
	return c.UploadSimple(url, nil, files...)
}
