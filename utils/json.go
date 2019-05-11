// Copyright 2019 Siva Chegondi

package utils

import (
	"encoding/json"
)

// Result json type to loop over keys
type Result struct {
	RequestID     string `json:"request_id"`
	LeaseID       string `json:"lease_id"`
	Renewable     bool
	LeaseDuration int `json:"lease_duration"`
	Data          data
	Metadata      resultMetadata
	WrapInfo      string `json:"wrap_info"`
	Warnings      string
	Auth          string
	Errors        []string
}

type data struct {
	Data     map[string]interface{}
	Metadata resultMetadata
	Keys     []string
}
type resultMetadata struct {
	CreatedTime  string `json:"created_time"`
	DeletionTime string `json:"deletion_time"`
	version      int
	Destroyed    bool
}

// Request json type to push data
type Request struct {
	Options  Options           `json:"options,omitempty"`
	Data     map[string]string `json:"data"`
	Versions []int             `json:"versions,omitempty"`
}

// Options check-and-set
type Options struct {
	Cas int `json:"cas,omitempty"`
}

// FormatResult unmarshals in to Result type
func FormatResult(data []byte) Result {
	var v Result
	if err := json.Unmarshal(data, &v); err != nil {
		panic(err)
	}
	return v
}

// CustomMarshal marshals based on type
func CustomMarshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}
