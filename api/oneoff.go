package api

import (
	"net/http"
	"time"
)

// Config represents some of the server's configuration.
type Config struct {
	// Auth holds information about how users can authenticate.
	//
	// NOTE: this can go one layer deeper after multi-auth.
	// https://github.com/scitran/core/pull/652
	Auth map[string]interface{} `json:"auth"`

	// Site holds multi-site registration information.
	// This feature is depreciated.
	Site map[string]interface{} `json:"site"`

	Created  time.Time `json:"created"`
	Modified time.Time `json:"modified"`
}

func (c *Client) GetConfig() (*Config, *http.Response, error) {
	var aerr *Error
	var config *Config
	resp, err := c.New().Get("config").Receive(&config, &aerr)
	return config, resp, Coalesce(err, aerr)
}

// Version identifies the upgrade level of system components.
type Version struct {
	// Database represents the database schema level.
	Database int `json:"database"`
}

func (c *Client) GetVersion() (*Version, *http.Response, error) {
	var aerr *Error
	var version *Version
	resp, err := c.New().Get("version").Receive(&version, &aerr)
	return version, resp, Coalesce(err, aerr)
}
