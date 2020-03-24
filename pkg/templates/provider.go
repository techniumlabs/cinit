package templates

import (
	"errors"

	"github.com/techniumlabs/cinit/pkg/templates/providers"
)

type TemplateProvider interface {
	ResolveTemplates(source string, dest string, vars map[string]string) (error)
}

type TemplateClient struct {
	Provider TemplateProvider
}

func NewTemplateClient(providerName string) (*TemplateClient, error) {

	if providerName == "" {
		providerName = "default"
	}

	switch providerName {
	case "default":
		return &TemplateClient{
			Provider: providers.NewDefaultTemplateProvider(),
		}, nil
	}

	return nil, errors.New("No Provider Found")
}
