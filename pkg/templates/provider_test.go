package templates

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewTemplatesClient(t *testing.T) {
	client, err := NewTemplateClient("default")
	assert.IsType(t, &TemplateClient{}, client, "Should be of client type")
	assert.NotNil(t, client.Provider, "Should not be nil")
	assert.NoError(t, err, "No Error")
}

func TestNewTemplatesClientWithEmptyConfig(t *testing.T) {
	client, err := NewTemplateClient("")
	assert.IsType(t, &TemplateClient{}, client, "Should be of client type")
	assert.NotNil(t, client.Provider, "Should not be nil")
	assert.NoError(t, err, "No Error")
}

func TestNewTemplatesClientWithInvalidConfig(t *testing.T) {
	client, err := NewTemplateClient("unknown-template")
	assert.Nil(t, client, "Should be nil")
	assert.EqualError(t, err, "No Provider Found", "Err should be same")
}
