package utils

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

func prepareRequest(vaultEndpoint, method string, data []byte) (*http.Client, *http.Request) {
	vaultKey := os.Getenv("CADDY_CLUSTERING_VAULT_KEY")
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
func LoadStore(vaultEndpoint string, data []byte) (Result, error) {
	client, req := prepareRequest(vaultEndpoint, "POST", data)

	res, err := client.Do(req)
	if err != nil {
		return Result{}, err
	}
	data, _ = ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	return FormatResult(data), nil
}

// DeleteStore deletes item with key
func DeleteStore(vaultEndpoint string) (Result, error) {
	client, req := prepareRequest(vaultEndpoint, "DELETE", nil)

	_, err := client.Do(req)
	return Result{}, err
}
