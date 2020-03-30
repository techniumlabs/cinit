package vault

import (
	"errors"
	"io/ioutil"
	"os"
	"strings"

	"fmt"

	vaultapi "github.com/hashicorp/vault/api"
	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/rest"
)

type VaultSecretProvider struct {
	Client *vaultapi.Client
}

const (
	serviceAccountFile = "/var/run/secrets/kubernetes.io/serviceaccount/token"
)

func NewVaultSecretProvider() (*VaultSecretProvider, error) {

	config := &vaultapi.Config{}

	// Check environment for VAULT_ADDR and VAULT_TOKEN
	envVaultAddr := os.Getenv("VAULT_ADDR")
	envVaultToken := os.Getenv("VAULT_TOKEN")

	if envVaultAddr != "" {
		config.Address = envVaultAddr
	} else {
		return nil, errors.New("VAULT_ADDR is null")
	}

	client, err := vaultapi.NewClient(config)
	if err != nil {
		log.Printf("Failed to create vault client. Error %s", err.Error())
		return nil, err
	}

	if envVaultToken != "" {
		log.Info("Using VAULT_TOKEN")
		client.SetToken(envVaultToken)
		return &VaultSecretProvider{
			Client: client,
		}, nil
	}

	homedir := os.Getenv("HOME")
	if _, err := os.Stat(homedir + "/.vault-token"); err == nil {
		tokendata, err := ioutil.ReadFile(homedir + "/.vault-token")
		if err != nil {
			log.Warnf("Could not read vault token from ~/.vault-token. %s", err.Error())
		} else {
			log.Info("Using ~/.vault-token")
			client.SetToken(string(tokendata))
			return &VaultSecretProvider{
				Client: client,
			}, nil
		}
	}

	// Check if the binary is running inside kuberentes container
	_, err = rest.InClusterConfig()
	if err != nil {
		return nil, err
	} else {
		err = VaultKubernetesAuth(client)
		if err != nil {
			return nil, fmt.Errorf("Failed to authenticate with K8s token %v", err)
		}
		return &VaultSecretProvider{
			Client: client,
		}, nil
	}

	return nil, errors.New("No auth possible")
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
			parsedString[key] = data[vaultKey].(string)
		} else {
			parsedString[key] = value
		}
	}
	parsedString["VAULT_TOKEN"] = v.Client.Token()
	return parsedString
}

func VaultKubernetesAuth(client *vaultapi.Client) error {
	jwt, err := ioutil.ReadFile(serviceAccountFile)
	if err != nil {
		return err
	}

	vaultRole := os.Getenv("VAULT_LOGIN_ROLE")
	if vaultRole == "" {
		return errors.New("VAULT_LOGIN_ROLE is null")
	}
	authPath := os.Getenv("KUBERNETES_AUTH_PATH")
	if authPath == "" {
		return errors.New("KUBERNETES_AUTH_PATH is null")
	}

	data := map[string]interface{}{"jwt": string(jwt), "role": vaultRole}

	secret, err := client.Logical().Write(fmt.Sprintf("auth/%s/login", authPath), data)
	if err != nil {
		log.Error("Failed to request new Vault token", err.Error())
		return err
	}

	if secret == nil {
		log.Error("received empty answer from Vault")
		return errors.New("Empty answer from vault")
	}

	client.SetToken(secret.Auth.ClientToken)
	return nil
}
