// Copyright 2019 Siva Chegondi

package utils

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

var req *http.Request

const vaultKey = "s.dNGuYevekYTJiUEupTfkwZhg"

// QueryStore connects to vault's pki store
func QueryStore(vaultEndpoint string) Result {
	var data []byte
	httpClient := &http.Client{
		Timeout: 30 * time.Second,
	}
	req, err := http.NewRequest("GET", vaultEndpoint, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Add("X-Vault-Token", vaultKey)

	if res, err := httpClient.Do(req); err != nil {
		panic(err)
	} else {
		data, _ = ioutil.ReadAll(res.Body)
		defer res.Body.Close()
	}
	return FormatResult(data)
}

// LoadStore loads store with data regarding
func LoadStore(vaultEndpoint string, data []byte) error {
	httpClient := &http.Client{}
	req, err := http.NewRequest("POST", vaultEndpoint, bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	req.Header.Add("X-Vault-Token", vaultKey)
	fmt.Println(req)

	res, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	data, _ = ioutil.ReadAll(res.Body)
	fmt.Print(string(data))
	defer res.Body.Close()
	return nil
}
