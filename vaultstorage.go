package caddyvault

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/caddyserver/caddy/caddytls"
	"github.com/mholt/certmagic"
	"github.com/siva-chegondi/caddyvault/utils"
)

const (
	defaultPrefix = "caddycerts"
	loadURL       = "/v1/" + defaultPrefix + "/data/"
	listURL       = "/v1/" + defaultPrefix + "/metadata/"
	storeURL      = "/v1/" + defaultPrefix + "/data/"
	deleteURL     = "/v1/" + defaultPrefix + "/metadata/"
)

// VaultStorage storage for ACME certificates
type VaultStorage struct {
	API string
}

func init() {
	caddytls.RegisterClusterPlugin("vault", constructVaultPlugin)
}

// creating instance of VaultStorage
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
		return list, certmagic.ErrNotExist(fmt.Errorf("List with prefix %s is empty", prefix))
	}
	return list, nil
}

// Load retrieves certificate of key
func (vaultStorage *VaultStorage) Load(key string) ([]byte, error) {
	res := utils.QueryStore(vaultStorage.API + loadURL + key)
	if len(res.Data.Data) == 0 {
		return []byte{}, certmagic.ErrNotExist(fmt.Errorf("Key %s does not exists", key))
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
		return errors.New("Failed to store, error: " + response.Errors[0])
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
	_, err := vaultStorage.List(key, false)
	_, ok := err.(certmagic.ErrNotExist)
	modified, merror := time.Parse(time.RFC3339, res.Data.Metadata.CreatedTime)
	return certmagic.KeyInfo{
		Key:        key,
		IsTerminal: !ok,
		Size:       int64(len(res.Data.Data[key].(string))),
		Modified:   modified,
	}, merror
}

/*
Util functions start here
listPath and queryPath
*/

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

// Lock locks operations on certificate with particular key
func (vaultStorage *VaultStorage) Lock(key string) error {
	key = key + ".lock"

	if vaultStorage.Exists(key) {

		if stat, err := vaultStorage.Stat(key); err == nil {

			// check for deadlock, wait for 5 (300s) minutes
			if time.Now().Unix()-stat.Modified.Unix() > 60 {
				vaultStorage.Unlock(key)
			} else {
				return errors.New("Lock already exists")
			}
		} else {
			return err
		}
	}

	return lockSystem(key, vaultStorage.API+storeURL+key)
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
		return errors.New("Failed to delete" + response.Errors[0])
	}
	return err
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
		return errors.New("Failed to lock: " + response.Errors[0])
	}
	return err
}
