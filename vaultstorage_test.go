package caddyvault_test

import (
	"os"
	"path"
	"testing"

	"github.com/siva-chegondi/caddyvault"
	"github.com/stretchr/testify/assert"
)

var vaultStorage *caddyvault.VaultStorage
var certPath string

func TestMain(m *testing.M) {
	os.Setenv("CADDY_CLUSTERING_VAULT_KEY", "s.1Dcdj2KeQbIbuibwGEhSBrQM")
	os.Setenv("CADDY_CLUSTERING_VAULT_ENDPOINT", "http://localhost:8200")
	vaultStorage = &caddyvault.VaultStorage{
		API: os.Getenv("CADDY_CLUSTERING_VAULT_ENDPOINT"),
	}
	certPath = path.Join("acme", "acme-v02.api.letsencrypt.org", "sites", "tls", "tls.crt")
	os.Exit(m.Run())
}

func TestStore(t *testing.T) {
	err := vaultStorage.Store(certPath, []byte("dump_data"))
	assert.NoError(t, err, "should store data")
}

func TestLoad(t *testing.T) {
	dataInBytes, _ := vaultStorage.Load(certPath)
	assert.Equal(t, `"dump_data"`, string(dataInBytes), "Did not found items")
}

func TestExists(t *testing.T) {
	status := vaultStorage.Exists(certPath)
	assert.True(t, status, "should exists")
}

func TestStat(t *testing.T) {
	keyInfo, err := vaultStorage.Stat(certPath)
	assert.NoError(t, err, "should not fail")
	assert.Equal(t, int64(len("dump_data")), keyInfo.Size, "key sizes should match")
}

func TestList(t *testing.T) {
	list, err := vaultStorage.List(certPath, false)
	assert.NoError(t, err, "should not fail listing")
	assert.NotEmpty(t, list, "list should not be empty")
}

func TestLock(t *testing.T) {
	err := vaultStorage.Lock(certPath)
	assert.NoError(t, err, "should not fail to lock")
	err = vaultStorage.Lock(certPath)
	assert.Error(t, err, "should fail to lock")
}

func TestUnlock(t *testing.T) {
	err := vaultStorage.Unlock(certPath)
	assert.NoError(t, err, "should not fail to unlock")
}

func TestDelete(t *testing.T) {
	err := vaultStorage.Delete(certPath)
	assert.NoError(t, err, "Should delete check")
}
