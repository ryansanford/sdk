package api

import (
	"time"
)

type File struct {
	Name   string  `json:"name,omitempty"`
	Origin *Origin `json:"origin,omitempty"`
	Size   int     `json:"size,omitempty"`

	Instrument   string   `json:"instrument,omitempty"`
	Mimetype     string   `json:"mimetype,omitempty"`
	Measurements []string `json:"measurements,omitempty"`
	Type         string   `json:"type,omitempty"`
	Tags         []string `json:"tags,omitempty"`

	Metadata map[string]interface{} `json:"metadata,omitempty"`

	Created  *time.Time `json:"created,omitempty"`
	Modified *time.Time `json:"modified,omitempty"`
}
