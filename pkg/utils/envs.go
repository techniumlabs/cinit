package utils

import "fmt"

func ConvertMapToEnv(envMap map[string]string) []string {
	var envs []string

	for key, value := range envMap {
		envs = append(envs, fmt.Sprintf("%s=%s", key, value))
	}

	return envs
}
