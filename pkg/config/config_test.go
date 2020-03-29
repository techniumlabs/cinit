package config

import (
	"io/ioutil"
	"testing"

	"os"

	"github.com/stretchr/testify/assert"
	"gotest.tools/fs"
)

func TestConfigFileInSameDirectory(t *testing.T) {
	dir := fs.NewDir(t, "test-config-file", fs.WithFile(".cinit.yaml", `
---
providers:
  - MyCustomProvider
`))
	defer dir.Remove()
	os.Chdir(dir.Path())
	c, _ := Load("")
	assert.Equal(t, c.SecretProviders[0], "MyCustomProvider", "Config should be from .cinit.yaml in current dir")
}

func TestConfigFileInHomeDirectory(t *testing.T) {
	home, _ := os.UserHomeDir()
	ioutil.WriteFile(home+"/.cinit.yaml", []byte(`
---
providers:
  - MyCustomProvider
`), 0644)
	defer os.RemoveAll(home + "/.cinit.yaml")
	c, _ := Load("")
	assert.Equal(t, c.SecretProviders[0], "MyCustomProvider", "Config should be from .cinit.yaml in home dir")
}

func TestConfigFileUserProvided(t *testing.T) {
	dir := fs.NewDir(t, "test-config-file", fs.WithFile(".cinit.yaml", `
---
providers:
  - MyCustomProvider
`))
	defer dir.Remove()
	cfgFile := dir.Path() + "/.cinit.yaml"
	c, _ := Load(cfgFile)
	assert.Equal(t, c.SecretProviders[0], "MyCustomProvider", "Config should be from provided file")
}

func TestConfigValue(t *testing.T) {
	dir := fs.NewDir(t, "test-config-file", fs.WithFile(".cinit.yaml", `
---
providers:
  - MyCustomProvider1
  - MyCustomProvider2
templates:
  - source: /tmp/src1
    dest: /tmp/dest1
  - source: /tmp/src2
    dest: /tmp/dest2
`))
	defer dir.Remove()
	cfgFile := dir.Path() + "/.cinit.yaml"
	c, _ := Load(cfgFile)
	assert.Equal(t, c.SecretProviders[0], "MyCustomProvider1", "Secret Provider should be equal")
	assert.Equal(t, c.SecretProviders[1], "MyCustomProvider2", "Secret Provider should be equal")
	assert.Equal(t, c.Templates[0].Source, "/tmp/src1", "Template source should be equal")
	assert.Equal(t, c.Templates[1].Source, "/tmp/src2", "Template source should be equal")
	assert.Equal(t, c.Templates[0].Dest, "/tmp/dest1", "Template dest should be equal")
	assert.Equal(t, c.Templates[1].Dest, "/tmp/dest2", "Template dest should be equal")
}
