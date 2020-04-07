package vault

import (
	"os"
	"testing"

	"io/ioutil"

	"gotest.tools/fs"

	vaultapi "github.com/hashicorp/vault/api"
	"github.com/stretchr/testify/assert"
)

func TestNewVaultProviderWithoutAddr(t *testing.T) {
	os.Setenv("VAULT_ADDR", "")

	provider, err := NewVaultSecretProvider()
	assert.Error(t, err, "Should return error for null VAULT_ADDR")
	assert.Nil(t, provider, "Should be nil")
}

func TestNewVaultProviderWithToken(t *testing.T) {
	token := "some-very-long-token"
	os.Setenv("VAULT_ADDR", "https://localhost:9090")
	os.Setenv("VAULT_TOKEN", token)

	provider, err := NewVaultSecretProvider()
	assert.NotNil(t, provider, "Should not be nil")
	assert.Nil(t, err, "Should be nil")
	assert.Equal(t, token, provider.Client.Token(), "Token should be equal")
}

func TestNewVaultProviderWithTokenFile(t *testing.T) {
	token := "some-very-long-token"
	os.Setenv("VAULT_ADDR", "https://localhost:9090")
	os.Setenv("VAULT_TOKEN", "")

	home, _ := os.UserHomeDir()
	ioutil.WriteFile(home+"/.vault-token", []byte(token), 0644)
	defer os.RemoveAll(home + "/.vault-token")

	provider, err := NewVaultSecretProvider()
	assert.NotNil(t, provider, "Should not be nil")
	assert.Nil(t, err, "Should be nil")
	assert.Equal(t, token, provider.Client.Token(), "Token should be equal")
}

func TestNewVaultProviderWithoutToken(t *testing.T) {
	os.Setenv("VAULT_ADDR", "https://localhost:9090")
	os.Setenv("VAULT_TOKEN", "")

	provider, err := NewVaultSecretProvider()
	assert.Error(t, err, "Should return error")
	assert.Nil(t, provider, "Should be nil")
}

func TestNewVaultProviderWithTokenAsPriority(t *testing.T) {
	token := "some-very-long-token"
	os.Setenv("VAULT_ADDR", "https://localhost:9090")
	os.Setenv("VAULT_TOKEN", token)

	home, _ := os.UserHomeDir()
	ioutil.WriteFile(home+"/.vault-token", []byte("some-other-token"), 0644)
	defer os.RemoveAll(home + "/.vault-token")

	provider, err := NewVaultSecretProvider()
	assert.NotNil(t, provider, "Should not be nil")
	assert.Nil(t, err, "Should be nil")
	assert.Equal(t, token, provider.Client.Token(), "Token should be equal")
}

func TestVaultKubernetesAuthWithoutServiceAccountFile(t *testing.T) {
	client := &vaultapi.Client{}
	err := VaultKubernetesAuth(client)
	assert.Error(t, err, "Should throw error on non existent serviceAccountFile")
}

func TestVaultKubernetesAuthWithoutRole(t *testing.T) {
	client := &vaultapi.Client{}
	dir := fs.NewDir(t, "serviceAccountFile", fs.WithFile("token", `some-token`))
	defer dir.Remove()

	serviceAccountFile = dir.Path() + "/token"
	ioutil.WriteFile(serviceAccountFile, []byte("some-token"), 0644)
	err := VaultKubernetesAuth(client)
	assert.EqualError(t, err, "VAULT_LOGIN_ROLE is null", "Should throw error on non existent Login Role")
}

func TestVaultKubernetesAuthWithoutPath(t *testing.T) {
	client := &vaultapi.Client{}
	dir := fs.NewDir(t, "serviceAccountFile", fs.WithFile("token", `some-token`))
	defer dir.Remove()

	serviceAccountFile = dir.Path() + "/token"
	ioutil.WriteFile(serviceAccountFile, []byte("some-token"), 0644)
	os.Setenv("VAULT_LOGIN_ROLE", "some-role")
	err := VaultKubernetesAuth(client)
	assert.EqualError(t, err, "KUBERNETES_AUTH_PATH is null", "Should throw error on non existent Login Role")
}
