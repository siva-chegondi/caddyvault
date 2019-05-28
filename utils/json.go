package utils

import (
	"encoding/json"
)

// Result json type to loop over keys
type Result struct {
	Renewable bool
	Warnings  string
	Auth      string
	Errors    []string

	Data     data
	Metadata resultMetadata

	LeaseID       string `json:"lease_id"`
	WrapInfo      string `json:"wrap_info"`
	RequestID     string `json:"request_id"`
	LeaseDuration int    `json:"lease_duration"`
}

type data struct {
	Keys     []string
	Metadata resultMetadata
	Data     map[string]interface{}
}
type resultMetadata struct {
	version      int
	Destroyed    bool
	CreatedTime  string `json:"created_time"`
	DeletionTime string `json:"deletion_time"`
}

// Request json type to push data
type Request struct {
	Data     map[string]string `json:"data"`
	Options  Options           `json:"options,omitempty"`
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
