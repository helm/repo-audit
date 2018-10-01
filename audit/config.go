package audit

import (
	"encoding/json"
	"io/ioutil"
)

// This file contains the configuration handling for the auditing tool. Config
// is passed in as JSON with an array of objects. The array enables multiple
// repositories to be checked and each object describes a repository.

// Configs contains the configuration for one or more repository
type Configs []*Config

// Config contains the configuration for a single repository
type Config struct {

	// The name of the repository. Useful when sending reports and identifying
	// a specific one
	Name string `json:"name"`

	// A URL to the repository
	Location string `json:"location"`
}

// LoadConfigFromFile loads configuration about repositories from a file on
// the local file system
func LoadConfigFromFile(path string) (Configs, error) {
	c, err := ioutil.ReadFile(path)
	if err != nil {
		// TODO(mattfarina): wrap the error into something more useful
		return nil, err
	}

	return LoadConfig(c)
}

// LoadConfig loads configuration about repositories from a byte stream
func LoadConfig(data []byte) (Configs, error) {
	cfg := Configs{}
	err := json.Unmarshal(data, &cfg)

	// TODO(mattfarina): wrap the error into something more useful if one exists
	return cfg, err
}
