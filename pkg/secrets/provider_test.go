package secrets

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/techniumlabs/cinit/pkg/config"
	"github.com/techniumlabs/cinit/pkg/secrets/providers/vault"
)

func TestNewSecretsClient(t *testing.T) {
	config := config.Config{
		SecretProviders: []string{"vault"},
	}

	client := NewSecretsClient(&config)
	assert.Len(t, client.Providers, 1, "Should have one provider")
	assert.IsType(t, &vault.VaultSecretProvider{}, client.Providers[0], "Should Create vault provider")
}

func TestNewSecretsClientWithEmptyConfig(t *testing.T) {
	config := config.Config{
		SecretProviders: []string{},
	}

	client := NewSecretsClient(&config)
	assert.Len(t, client.Providers, 1, "Should have one provider")
	assert.IsType(t, &vault.VaultSecretProvider{}, client.Providers[0], "Should Create vault provider")
}

func TestNewSecretsClientWithInvalidProvider(t *testing.T) {
	config := config.Config{
		SecretProviders: []string{"someprovider"},
	}

	client := NewSecretsClient(&config)
	assert.Len(t, client.Providers, 0, "Should have no provider")
}

func TestNewSecretsClientWithBothValidAndInvalid(t *testing.T) {
	config := config.Config{
		SecretProviders: []string{"someprovider", "vault"},
	}

	client := NewSecretsClient(&config)
	assert.Len(t, client.Providers, 1, "Should have one provider")
	assert.IsType(t, &vault.VaultSecretProvider{}, client.Providers[0], "Should Create vault provider")
}
