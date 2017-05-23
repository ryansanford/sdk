package api

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// ParseApiKey accepts an API key and returns the hostname, port, key, and any parsing error.
func ParseApiKey(apiKey string) (string, int, string, error) {
	var err error
	host := ""
	port := 443
	key := ""

	splits := strings.Split(apiKey, ":")

	if len(splits) < 2 {
		return host, port, key, errors.New("Invalid API key")
	}

	if len(splits) == 2 {
		host = splits[0]
		key = splits[1]
	} else {
		host = splits[0]
		port, err = strconv.Atoi(splits[1])
		key = splits[len(splits)-1]
	}

	return host, port, key, err
}

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
