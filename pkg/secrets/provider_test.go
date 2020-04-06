package secrets

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/techniumlabs/cinit/pkg/config"
	"github.com/techniumlabs/cinit/pkg/secrets/providers/vault"
)

func TestNewSecretsClient(t *testing.T) {
	config := config.Config{
		SecretProviders: []string{"vault"},
	}

	os.Setenv("VAULT_ADDR", "https://localhost:9090")
	os.Setenv("VAULT_TOKEN", "some-token")
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
	assert.Nil(t, client, "Invalid Provider")
}

func TestNewSecretsClientWithBothValidAndInvalid(t *testing.T) {
	config := config.Config{
		SecretProviders: []string{"someprovider", "vault"},
	}

	client := NewSecretsClient(&config)
	assert.Nil(t, client, "Invalid Provider")
}

type MockSecretProvider struct {
	mock.Mock
}

func (m *MockSecretProvider) ResolveSecrets(envMap map[string]string) map[string]string {
	args := m.Called(envMap)
	return args.Get(0).(map[string]string)
}

func TestEnvParsing(t *testing.T) {
	config := config.Config{
		SecretProviders: []string{"vault"},
	}

	client := NewSecretsClient(&config)
	assert.Len(t, client.Providers, 1, "Should have one provider")

	mockProvider := new(MockSecretProvider)
	os.Clearenv()
	envMap := map[string]string{"var1": "value1"}
	os.Setenv("var1", "value1")
	mockProvider.On("ResolveSecrets", envMap).Return(envMap)

	client.Providers[0] = mockProvider
	client.GetParsedEnvs()
	mockProvider.AssertNumberOfCalls(t, "ResolveSecrets", 1)

}
