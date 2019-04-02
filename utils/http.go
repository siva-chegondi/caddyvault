// Copyright 2019 Siva Chegondi

package utils

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

func prepareRequest(vaultEndpoint, method string, data []byte) (*http.Client, *http.Request) {
	const vaultKey = "s.5B99YC7ZX9yhDZ7wxpgvFvtN"
	httpClient := &http.Client{
		Timeout: 30 * time.Second,
	}
	req, err := http.NewRequest(method, vaultEndpoint, bytes.NewBuffer(data))
	if err != nil {
		panic(err)
	}
	req.Header.Add("X-Vault-Token", vaultKey)
	return httpClient, req
}

// QueryStore connects to vault's pki store
func QueryStore(vaultEndpoint string) Result {
	var data []byte
	client, req := prepareRequest(vaultEndpoint, "GET", nil)
	if res, err := client.Do(req); err != nil {
		panic(err)
	} else {
		data, _ = ioutil.ReadAll(res.Body)
		defer res.Body.Close()
	}
	return FormatResult(data)
}

// ListStore connects to vault's pki store
func ListStore(vaultEndpoint string) Result {
	var data []byte
	client, req := prepareRequest(vaultEndpoint, "LIST", nil)

	if res, err := client.Do(req); err != nil {
		panic(err)
	} else {
		data, _ = ioutil.ReadAll(res.Body)
		defer res.Body.Close()
	}
	return FormatResult(data)
}

// LoadStore loads store with data regarding
func LoadStore(vaultEndpoint string, data []byte) error {
	client, req := prepareRequest(vaultEndpoint, "POST", data)

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	data, _ = ioutil.ReadAll(res.Body)
	fmt.Print(string(data))
	defer res.Body.Close()
	return nil
}

// DeleteStore deletes item with key
func DeleteStore(vaultEndpoint string, data []byte) error {
	client, req := prepareRequest(vaultEndpoint, "DELETE", data)

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	return nil
}
