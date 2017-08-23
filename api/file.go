package api

import (
	"time"
)

type File struct {
	Name   string  `json:"name,omitempty"`
	Origin *Origin `json:"origin,omitempty"`
	Size   int     `json:"size,omitempty"`

	Modality     string   `json:"modality,omitempty"`
	Mimetype     string   `json:"mimetype,omitempty"`
	Measurements []string `json:"measurements,omitempty"`
	Type         string   `json:"type,omitempty"`
	Tags         []string `json:"tags,omitempty"`

	Info map[string]interface{} `json:"info,omitempty"`

	Created  *time.Time `json:"created,omitempty"`
	Modified *time.Time `json:"modified,omitempty"`
}

type FileFields struct {
	Modality     string   `json:"modality,omitempty"`
	Measurements []string `json:"measurements,omitempty"`
	Type         string   `json:"type,omitempty"`
}
