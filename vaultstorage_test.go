package caddyvault_test

import (
	"os"
	"path"
	"testing"

	"github.com/siva-chegondi/caddyvault"
	"github.com/stretchr/testify/assert"
)

var (
	certPath     string
	vaultStorage *caddyvault.VaultStorage
)

const certData = `
-----BEGIN PUBLIC KEY-----
MIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEAlRuRnThUjU8/prwYxbty
WPT9pURI3lbsKMiB6Fn/VHOKE13p4D8xgOCADpdRagdT6n4etr9atzDKUSvpMtR3
CP5noNc97WiNCggBjVWhs7szEe8ugyqF23XwpHQ6uV1LKH50m92MbOWfCtjU9p/x
qhNpQQ1AZhqNy5Gevap5k8XzRmjSldNAFZMY7Yv3Gi+nyCwGwpVtBUwhuLzgNFK/
yDtw2WcWmUU7NuC8Q6MWvPebxVtCfVp/iQU6q60yyt6aGOBkhAX0LpKAEhKidixY
nP9PNVBvxgu3XZ4P36gZV6+ummKdBVnc3NqwBLu5+CcdRdusmHPHd5pHf4/38Z3/
6qU2a/fPvWzceVTEgZ47QjFMTCTmCwNt29cvi7zZeQzjtwQgn4ipN9NibRH/Ax/q
TbIzHfrJ1xa2RteWSdFjwtxi9C20HUkjXSeI4YlzQMH0fPX6KCE7aVePTOnB69I/
a9/q96DiXZajwlpq3wFctrs1oXqBp5DVrCIj8hU2wNgB7LtQ1mCtsYz//heai0K9
PhE4X6hiE0YmeAZjR0uHl8M/5aW9xCoJ72+12kKpWAa0SFRWLy6FejNYCYpkupVJ
yecLk/4L1W0l6jQQZnWErXZYe0PNFcmwGXy1Rep83kfBRNKRy5tvocalLlwXLdUk
AIU+2GKjyT3iMuzZxxFxPFMCAwEAAQ==
-----END PUBLIC KEY-----
`

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
	err := vaultStorage.Store(certPath, []byte(certData))
	assert.NoError(t, err, "should store data")
}

func TestLoad(t *testing.T) {
	dataInBytes, _ := vaultStorage.Load(certPath)
	assert.Equal(t, certData, string(dataInBytes), "Did not found items")
}

func TestExists(t *testing.T) {
	status := vaultStorage.Exists(certPath)
	assert.True(t, status, "should exists")
}

func TestStat(t *testing.T) {
	keyInfo, err := vaultStorage.Stat(certPath)
	assert.NoError(t, err, "should not fail")
	assert.Equal(t, int64(len(certData)), keyInfo.Size, "key sizes should match")
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
