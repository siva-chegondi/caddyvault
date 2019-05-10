// Copyright 2019 Siva Chegondi
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package caddyvault

import (
	"encoding/json"
	"errors"
	"os"
	"strings"
	"time"

	"github.com/mholt/caddy/caddytls"
	"github.com/mholt/certmagic"
	"github.com/siva-chegondi/caddyvault/utils"
)

const (
	loadURL   = "/v1/caddycerts/data/"
	listURL   = "/v1/caddycerts/metadata/"
	storeURL  = "/v1/caddycerts/data/"
	deleteURL = "/v1/caddycerts/metadata/"
)

// VaultStorage storage for ACME certificates
type VaultStorage struct {
	API string
}

func init() {
	caddytls.RegisterClusterPlugin("vault", constructVaultPlugin)
}

func constructVaultPlugin() (certmagic.Storage, error) {
	return &VaultStorage{
		API: os.Getenv("CADDY_CLUSTERING_VAULT_ENDPOINT"),
	}, nil
}

// List lists certificates
func (vaultStorage *VaultStorage) List(prefix string, recursive bool) ([]string, error) {
	var list []string
	if recursive {
		list = listPath(vaultStorage.API+listURL, vaultStorage.API+loadURL, prefix)
	} else {
		list = queryPath(vaultStorage.API+loadURL, prefix)
	}

	if len(list) == 0 {
		return list, os.ErrNotExist
	}
	return list, nil
}

// Load retrieves certificate of key
func (vaultStorage *VaultStorage) Load(key string) ([]byte, error) {
	res := utils.QueryStore(vaultStorage.API + loadURL + key)
	if len(res.Data.Data) == 0 {
		return []byte{}, os.ErrNotExist
	}
	return []byte(res.Data.Data[key].(string)), nil
}

// Store stores certificate with key association
func (vaultStorage *VaultStorage) Store(key string, value []byte) error {
	data := make(map[string]string)
	data[key] = string(value)
	req := &utils.Request{
		Data: data,
	}
	byteData, _ := json.Marshal(req)
	response, err := utils.LoadStore(vaultStorage.API+storeURL+key, byteData)
	if len(response.Errors) > 0 {
		return errors.New(response.Errors[0])
	}
	return err
}

// Exists returns existance of certificate with key
func (vaultStorage *VaultStorage) Exists(key string) bool {
	res := utils.QueryStore(vaultStorage.API + loadURL + key)
	return len(res.Data.Data) > 0 && !res.Data.Metadata.Destroyed
}

// Stat retrieves status of certificate with key param
func (vaultStorage *VaultStorage) Stat(key string) (certmagic.KeyInfo, error) {
	res := utils.QueryStore(vaultStorage.API + loadURL + key)
	list, err := vaultStorage.List(key, false)
	modified, err := time.Parse(time.RFC3339, res.Data.Metadata.CreatedTime)
	return certmagic.KeyInfo{
		Key:        key,
		IsTerminal: len(list) > 0,
		Size:       int64(len(res.Data.Data[key].(string))),
		Modified:   modified,
	}, err
}

// Lock locks operations on certificate with particular key
func (vaultStorage *VaultStorage) Lock(key string) error {
	// check for deadlock, wait for 5 (300s) minutes
	key = key + ".lock"
	if vaultStorage.Exists(key) {
		if stat, err := vaultStorage.Stat(key); err == nil {
			if time.Now().Unix()-stat.Modified.Unix() > 300 {
				vaultStorage.Unlock(key)
			} else {
				return errors.New("Lock already exists")
			}
		} else {
			return err
		}
	}

	return lockSystem(key, vaultStorage.API+loadURL+key)
}

// Unlock unlocks operations on certificate data
func (vaultStorage *VaultStorage) Unlock(key string) error {
	if strings.Index(key, ".lock") < 0 {
		key = key + ".lock"
	}
	return vaultStorage.Delete(key)
}

// Delete deletes the certificate from vault.
func (vaultStorage *VaultStorage) Delete(key string) error {
	response, err := utils.DeleteStore(vaultStorage.API + deleteURL + key)
	if len(response.Errors) > 0 {
		return errors.New(response.Errors[0])
	}
	return err
}

func listPath(listurl, loadurl, prefix string) []string {
	var list []string
	var res utils.Result

	// list all the keys
	list = append(list, queryPath(loadurl, prefix)...)

	// list all the paths and loop keys
	res = utils.ListStore(listurl + prefix)
	for _, path := range res.Data.Keys {
		list = append(list, listPath(listurl+prefix, loadurl+prefix, "/"+path)...)
	}
	return list
}

func queryPath(url, prefix string) []string {
	var res utils.Result
	var list []string
	res = utils.QueryStore(url + prefix)
	for item := range res.Data.Data {
		list = append(list, item)
	}
	return list
}

func lockSystem(key, lockPath string) error {
	data := make(map[string]string)
	data[key] = "locked"
	postBody := utils.Request{
		Options: utils.Options{
			Cas: 0,
		},
		Data: data,
	}
	jsonData, _ := json.Marshal(postBody)
	response, err := utils.LoadStore(lockPath, jsonData)
	if len(response.Errors) > 0 {
		return errors.New(response.Errors[0])
	}
	return err
}
