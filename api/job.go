package api

import (
	"net/http"
	"time"
)

// Formula describes a single unit of work.
type Formula struct {
	Inputs  []*Input  `json:"inputs"`
	Target  Target    `json:"target"`
	Outputs []*Output `json:"outputs"`
}

// Input describes an asset that must be present for the forumla to execute.
type Input struct {
	Type     string `json:"type,omitempty"`
	URI      string `json:"uri,omitempty"`
	Location string `json:"location,omitempty"`
	VuID     string `json:"vu,omitempty"`
}

// Target describes what the formula will execute.
type Target struct {
	Command []string          `json:"command,omitempty"`
	Env     map[string]string `json:"env,omitempty"`
	Dir     string            `json:"dir,omitempty"`
}

// Output describes the creation of an asset after a formula is complete.
type Output struct {
	Type     string `json:"type,omitempty"`
	URI      string `json:"uri,omitempty"`
	Location string `json:"location,omitempty"`
	VuID     string `json:"vu,omitempty"`
}

// Result describes the result of a formula.
type Result struct {
	ExitCode int `json:"exitcode"`
}

// FormulaResult combines a (possibly-modified) Formula with any results.
type FormulaResult struct {
	Formula

	Result Result `json:"result"`
}

// Enum for job states.
type JobState string

const (
	Pending   JobState = "pending"
	Running   JobState = "running"
	Failed    JobState = "failed"
	Complete  JobState = "complete"
	Cancelled JobState = "cancelled"
)

// Enum for job retrieval attempts.
type JobRetrieval int

const (
	JobAquired JobRetrieval = iota
	NoPendingJobs
	JobFailure
)

type Job struct {
	Id     string `json:"id,omitempty"`
	GearId string `json:"gear_id,omitempty"`

	State   JobState `json:"state,omitempty"`
	Attempt int      `json:"attempt,omitempty"`
	Origin  *Origin  `json:"origin,omitempty"`

	Config      map[string]interface{} `json:"config,omitempty"`
	Inputs      map[string]interface{} `json:"inputs,omitempty"`
	Destination *ContainerReference    `json:"destination,omitempty"`
	Tags        []string               `json:"tags,omitempty"`

	Request *Formula `json:"request,omitempty"`

	ResultMetadata map[string]interface{} `json:"produced_metadata,omitempty"`
	ResultFiles    []string               `json:"saved_files,omitempty"`

	Created  *time.Time `json:"created,omitempty"`
	Modified *time.Time `json:"modified,omitempty"`
}

// Get all jobs endpoint is not implemented as it returns a different format
// https://github.com/scitran/core/issues/704

func (c *Client) GetJob(id string) (*Job, *http.Response, error) {
	var aerr *Error
	var job *Job

	// Should not require root flag
	// https://github.com/scitran/core/issues/657
	resp, err := c.New().Get("jobs/"+id+"?root=true").Receive(&job, &aerr)
	return job, resp, Coalesce(err, aerr)
}

func (c *Client) AddJob(job *Job) (string, *http.Response, error) {
	var aerr *Error
	var response *IdResponse
	var result string

	resp, err := c.New().Post("jobs/add").BodyJSON(job).Receive(&response, &aerr)

	if response != nil {
		result = response.Id
	}

	return result, resp, Coalesce(err, aerr)
}

func (c *Client) ModifyJob(id string, job *Job, asRoot bool) (*http.Response, error) {
	var aerr *Error

	// Should not have to specify root flag
	// https://github.com/scitran/core/issues/657
	url := "jobs/" + id
	if asRoot {
		url += "?root=true"
	}

	resp, err := c.New().Put(url).BodyJSON(job).Receive(nil, &aerr)
	return resp, Coalesce(err, aerr)
}

func (c *Client) StartNextPendingJob(tags ...string) (JobRetrieval, *Job, *http.Response, error) {
	var aerr *Error
	var job *Job

	params := &struct {
		Tags []string `url:"tags,omitempty"`
	}{
		Tags: tags,
	}

	// Should not require root flag
	// https://github.com/scitran/core/issues/657
	resp, err := c.New().Get("jobs/next"+"?root=true").QueryStruct(params).Receive(&job, &aerr)
	rerr := Coalesce(err, aerr)

	if rerr == nil && job != nil {
		return JobAquired, job, resp, nil
	} else if rerr != nil && resp.StatusCode == 400 {
		return NoPendingJobs, nil, resp, nil
	} else {
		return JobFailure, nil, resp, rerr
	}
}

func (c *Client) HeartbeatJob(id string) (*http.Response, error) {
	var aerr *Error

	// Send empty modification
	empty := map[string]string{}

	// Should not require root flag
	// https://github.com/scitran/core/issues/657
	resp, err := c.New().Put("jobs/"+id+"?root=true").BodyJSON(empty).Receive(nil, &aerr)
	return resp, Coalesce(err, aerr)
}

func (c *Client) ChangeJobState(id string, state JobState, asRoot bool) (*http.Response, error) {
	var aerr *Error
	jobMod := &Job{
		State: state,
	}

	// Should not have to specify root flag
	// https://github.com/scitran/core/issues/657
	url := "jobs/" + id
	if asRoot {
		url += "?root=true"
	}

	resp, err := c.New().Put(url).BodyJSON(jobMod).Receive(nil, &aerr)
	return resp, Coalesce(err, aerr)
}
