package secrets

import (
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/techniumlabs/cinit/pkg/config"
	"github.com/techniumlabs/cinit/pkg/secrets/providers/vault"
)

// SecretsProvider secrets provider interface
type SecretsProvider interface {
	ResolveSecrets(vars map[string]string) map[string]string
}

type SecretsClient struct {
	Providers []SecretsProvider
}

func NewSecretsClient(config *config.Config) *SecretsClient {
	client := &SecretsClient{}
	err := client.InitProviders(config.SecretProviders)
	if err != nil {
		log.Printf("Could not Initialize Providers %v", err)
	}

	return client
}

func (c *SecretsClient) InitProviders(providerNames []string) error {
	var providers []SecretsProvider
	if len(providerNames) == 0 {
		providerNames = []string{"vault"}
	}
	for _, providerName := range providerNames {
		if providerName == "vault" {
			provider, err := vault.NewVaultSecretProvider()
			if err != nil {
				log.Printf("%s", err.Error())
			} else {
				providers = append(providers, provider)
			}
		}
	}

	c.Providers = providers

	return nil
}

func (c *SecretsClient) GetParsedEnvs() map[string]string {
	var parsedMap map[string]string

	envMap := make(map[string]string)
	// Parse the env vars
	envs := os.Environ()
	for _, env := range envs {
		kv := strings.Split(env, "=")
		key, value := kv[0], kv[1]
		envMap[key] = value
	}

	for _, provider := range c.Providers {
		parsedMap = provider.ResolveSecrets(envMap)
	}

	return parsedMap
}
