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

package vaultstorage

import (
	"fmt"

	"github.com/mholt/certmagic"
	"github.com/siva-chegondi/caddyvault/utils"
)

const (
	loadURL  = "/v1/secret/data/"
	listURL  = "/v1/secret/metadata/"
	storeURL = "/v1/secret/data/"
)

// VaultStorage storage for ACME certificates
type VaultStorage struct {
	API string
}

// List lists certificates
func (vaultStorage *VaultStorage) List(prefix string, recursive bool) ([]string, error) {
	data := utils.QueryStore(vaultStorage.API + listURL + prefix)
	fmt.Println(data)
	return nil, nil
}

// Load retrieves certificate of key
func (vaultStorage *VaultStorage) Load(key string) ([]byte, error) {
	data := utils.QueryStore(vaultStorage.API + loadURL + key)
	return data.Data.Data["test"].([]byte), nil
}

// Store stores certificate with key association
func (vaultStorage *VaultStorage) Store(key string, value []byte) error {
	return utils.LoadStore(vaultStorage.API+storeURL+key, []byte(`{"data":{"cert":"`+string(value)+`"}}`))
}

// Exists returns existance of certificate with key
func (vaultStorage *VaultStorage) Exists(key string) bool {
	data := utils.QueryStore(vaultStorage.API + loadURL + key)
	return data.Metadata.Destroyed
}

// Stat retrieves status of certificate with key param
func (vaultStorage *VaultStorage) Stat(key string) (certmagic.KeyInfo, error) {
	data := utils.QueryStore(vaultStorage.API + loadURL + key)
	return certmagic.KeyInfo{
		Key:        key,
		IsTerminal: false,
		Size:       int64(len(data.Data.Data["test"].(string))),
		Modified:   data.Metadata.Created_time,
	}, nil
}

// // Lock locks operations on certificate with particular key
// func (vaultStorage *VaultStorage) Lock(key string) error {
// 	return nil
// }

// // Unlock unlocks operations on certificate data
// func (vaultStorage *VaultStorage) Unlock(key string) error {
// 	return nil
// }
