package api

import (
	"errors"
	"net/http"
	"strconv"
	"time"
)

type Collection struct {
	Id          string `json:"_id,omitempty"`
	Name        string `json:"label,omitempty"`
	Curator     string `json:"curator,omitempty"`
	Description string `json:"description,omitempty"`

	Created  *time.Time `json:"created,omitempty"`
	Modified *time.Time `json:"modified,omitempty"`
	Files    []*File    `json:"files,omitempty"`

	Notes []*Note                `json:"notes,omitempty"`
	Info  map[string]interface{} `json:"info,omitempty"`

	Public      *bool         `json:"public,omitempty"`
	Archived    *bool         `json:"archived,omitempty"`
	Permissions []*Permission `json:"permissions,omitempty"`
}

func (c *Client) GetAllCollections() ([]*Collection, *http.Response, error) {
	var aerr *Error
	var collections []*Collection
	resp, err := c.New().Get("collections").Receive(&collections, &aerr)
	return collections, resp, Coalesce(err, aerr)
}

func (c *Client) GetCollection(id string) (*Collection, *http.Response, error) {
	var aerr *Error
	var collection *Collection
	resp, err := c.New().Get("collections/"+id).Receive(&collection, &aerr)
	return collection, resp, Coalesce(err, aerr)
}

func (c *Client) GetCollectionSessions(id string) ([]*Session, *http.Response, error) {
	var aerr *Error
	var sessions []*Session
	resp, err := c.New().Get("collections/"+id+"/sessions").Receive(&sessions, &aerr)
	return sessions, resp, Coalesce(err, aerr)
}

func (c *Client) GetCollectionAcquisitions(id string) ([]*Session, *http.Response, error) {
	var aerr *Error
	var sessions []*Session
	resp, err := c.New().Get("collections/"+id+"/acquisitions").Receive(&sessions, &aerr)
	return sessions, resp, Coalesce(err, aerr)
}

func (c *Client) GetCollectionSessionAcquisitions(id string, sid string) ([]*Session, *http.Response, error) {
	var aerr *Error
	var sessions []*Session
	resp, err := c.New().Get("collections/"+id+"/acquisitions?session="+sid).Receive(&sessions, &aerr)
	return sessions, resp, Coalesce(err, aerr)
}

func (c *Client) AddCollection(collection *Collection) (string, *http.Response, error) {
	var aerr *Error
	var response *IdResponse
	var result string

	resp, err := c.New().Post("collections").BodyJSON(collection).Receive(&response, &aerr)

	if response != nil {
		result = response.Id
	}

	return result, resp, Coalesce(err, aerr)
}

func (c *Client) AddCollectionNote(id, text string) (*http.Response, error) {
	var aerr *Error
	var response *ModifiedResponse

	note := &Note{
		Text: text,
	}

	resp, err := c.New().Post("collections/"+id+"/notes").BodyJSON(note).Receive(&response, &aerr)

	// Should not have to check this count
	// https://github.com/scitran/core/issues/680
	if err == nil && aerr == nil && response.ModifiedCount != 1 {
		return resp, errors.New("Modifying collection " + id + " returned " + strconv.Itoa(response.ModifiedCount) + " instead of 1")
	}

	return resp, Coalesce(err, aerr)
}

func (c *Client) AddCollectionTag(id, tag string) (*http.Response, error) {
	var aerr *Error
	var response *ModifiedResponse

	var tagDoc interface{}
	tagDoc = map[string]interface{}{
		"value": tag,
	}

	resp, err := c.New().Post("collections/"+id+"/tags").BodyJSON(tagDoc).Receive(&response, &aerr)

	// Should not have to check this count
	// https://github.com/scitran/core/issues/680
	if err == nil && aerr == nil && response.ModifiedCount != 1 {
		return resp, errors.New("Modifying collection " + id + " returned " + strconv.Itoa(response.ModifiedCount) + " instead of 1")
	}

	return resp, Coalesce(err, aerr)
}

func (c *Client) ModifyCollection(id string, collection *Collection) (*http.Response, error) {
	var aerr *Error
	var response *ModifiedResponse

	resp, err := c.New().Put("collections/"+id).BodyJSON(collection).Receive(&response, &aerr)

	// Should not have to check this count
	// https://github.com/scitran/core/issues/680
	if err == nil && aerr == nil && response.ModifiedCount != 1 {
		return resp, errors.New("Modifying collection " + id + " returned " + strconv.Itoa(response.ModifiedCount) + " instead of 1")
	}

	return resp, Coalesce(err, aerr)
}

func (c *Client) DeleteCollection(id string) (*http.Response, error) {
	var aerr *Error

	// Unlike other delete endpoints, this doesn't return anything. Which is good, but inconsistent.
	// https://github.com/scitran/core/issues/680
	resp, err := c.New().Delete("collections/"+id).Receive(nil, &aerr)

	return resp, Coalesce(err, aerr)
}

func (c *Client) UploadToCollection(id string, files ...*UploadSource) (chan int64, chan error) {
	url := "collections/" + id + "/files"
	return c.UploadSimple(url, nil, files...)
}

func (c *Client) DownloadFromCollection(id string, filename string, destination *DownloadSource) (chan int64, chan error) {
	url := "collections/" + id + "/files/" + filename
	return c.DownloadSimple(url, destination)
}

// No progress reporting
func (c *Client) UploadFileToCollection(id string, path string) error {
	src := CreateUploadSourceFromFilenames(path)
	progress, result := c.UploadToCollection(id, src...)

	// drain and report
	for range progress {
	}
	return <-result
}

// No progress reporting
func (c *Client) DownloadFileFromCollection(id, name string, path string) error {
	src := CreateDownloadSourceFromFilename(path)
	progress, result := c.DownloadFromCollection(id, name, src)

	// drain and report
	for range progress {
	}
	return <-result
}
