package vault

import (
	"os"
	"testing"

	"io/ioutil"

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
