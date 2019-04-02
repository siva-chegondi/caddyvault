// Copyright 2019 Siva Chegondi

package utils

import (
	"encoding/json"
	"time"
)

// Result json type to loop over keys
type Result struct {
	Request_id     string
	Lease_id       string
	Renewable      bool
	Lease_duration int
	Data           map[string]interface{}
	Metadata       resultMetadata
	Wrap_info      string
	Warnings       string
	Auth           string
}

type resultMetadata struct {
	Created_time  time.Time
	Deletion_time time.Time
	version       int
	Destroyed     bool
}

// FormatResult unmarshals in to Result type
func FormatResult(data []byte) Result {
	var v Result
	if err := json.Unmarshal(data, &v); err != nil {
		panic(err)
	}
	return v
}
