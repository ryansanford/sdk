package api

import (
	"errors"
	"net/http"
	"strconv"
	"time"
)

type Acquisition struct {
	Id        string `json:"_id,omitempty"`
	Name      string `json:"label,omitempty"`
	SessionId string `json:"session,omitempty"`

	Timestamp *time.Time `json:"timestamp,omitempty"`

	Created  *time.Time `json:"created,omitempty"`
	Modified *time.Time `json:"modified,omitempty"`
	Files    []*File    `json:"files,omitempty"`

	Public      bool          `json:"public,omitempty"`
	Permissions []*Permission `json:"permissions,omitempty"`
}

func (c *Client) GetAllAcquisitions() ([]*Acquisition, *http.Response, error) {
	var aerr *Error
	var acquisitions []*Acquisition
	resp, err := c.New().Get("acquisitions").Receive(&acquisitions, &aerr)
	return acquisitions, resp, Coalesce(err, aerr)
}

func (c *Client) GetAcquisition(id string) (*Acquisition, *http.Response, error) {
	var aerr *Error
	var acquisition *Acquisition
	resp, err := c.New().Get("acquisitions/"+id).Receive(&acquisition, &aerr)
	return acquisition, resp, Coalesce(err, aerr)
}

func (c *Client) AddAcquisition(acquisition *Acquisition) (string, *http.Response, error) {
	var aerr *Error
	var response *IdResponse
	var result string

	resp, err := c.New().Post("acquisitions").BodyJSON(acquisition).Receive(&response, &aerr)

	if response != nil {
		result = response.Id
	}

	return result, resp, Coalesce(err, aerr)
}

func (c *Client) ModifyAcquisition(id string, acquisition *Acquisition) (*http.Response, error) {
	var aerr *Error
	var response *ModifiedResponse

	resp, err := c.New().Put("acquisitions/"+id).BodyJSON(acquisition).Receive(&response, &aerr)

	// Should not have to check this count
	// https://github.com/scitran/core/issues/680
	if err == nil && aerr == nil && response.ModifiedCount != 1 {
		return resp, errors.New("Modifying acquisition " + id + " returned " + strconv.Itoa(response.ModifiedCount) + " instead of 1")
	}

	return resp, Coalesce(err, aerr)
}

func (c *Client) DeleteAcquisition(id string) (*http.Response, error) {
	var aerr *Error
	var response *DeletedResponse

	resp, err := c.New().Delete("acquisitions/"+id).Receive(&response, &aerr)

	// Should not have to check this count
	// https://github.com/scitran/core/issues/680
	if err == nil && aerr == nil && response.DeletedCount != 1 {
		return resp, errors.New("Deleting acquisition " + id + " returned " + strconv.Itoa(response.DeletedCount) + " instead of 1")
	}

	return resp, Coalesce(err, aerr)
}

func (c *Client) UploadToAcquisition(id string, files ...*UploadSource) (chan int64, chan error) {
	url := "acquisitions/" + id + "/files"
	return c.UploadSimple(url, nil, files...)
}
