package vault

import (
	"log"
	"os"
	"strings"

	vaultapi "github.com/hashicorp/vault/api"
)

type VaultSecretProvider struct {
	Client *vaultapi.Client
}

func NewVaultSecretProvider() (*VaultSecretProvider, error) {

	config := &vaultapi.Config{}

	// Check environment for VAULT_ADDR and VAULT_TOKEN
	envVaultAddr := os.Getenv("VAULT_ADDR")
	envVaultToken := os.Getenv("VAULT_TOKEN")

	if envVaultAddr != "" {
		config.Address = envVaultAddr
	}

	client, err := vaultapi.NewClient(config)
	if err != nil {
		log.Printf("Failed to create vault client. Error %s", err.Error())
		return nil, err
	}

	if envVaultToken != "" {
		client.SetToken(envVaultToken)
	}

	return &VaultSecretProvider{
		Client: client,
	}, nil
}

func (v *VaultSecretProvider) ResolveSecrets(vars map[string]string) map[string]string {

	parsedString := make(map[string]string)
	// Fetch all Keys
	for key, value := range vars {
		if strings.HasPrefix(value, "vault:") {
			vaultKeyPath := strings.TrimPrefix(value, "vault:")
			vaultKeyPathArr := strings.Split(vaultKeyPath, "?")
			vaultPath, vaultKey := vaultKeyPathArr[0], vaultKeyPathArr[1]
			secret, err := v.Client.Logical().Read(vaultPath)
			if err != nil {
				log.Printf("Could not resolv %s. Err %s", vaultPath, err)
				parsedString[key] = value
			}

			data := secret.Data["data"].(map[string]interface{})
			parsedString[vaultKey] = data[vaultKey].(string)
		} else {
			parsedString[key] = value
		}
	}
	return parsedString
}
