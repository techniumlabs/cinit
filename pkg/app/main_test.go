package app

import (
	"os"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"gotest.tools/fs"
)

func TestNewApp(t *testing.T) {
	dir := fs.NewDir(t, "test-config-file", fs.WithFile(".cinit.yaml", ``))
	defer dir.Remove()
	os.Chdir(dir.Path())
	app, err := NewApp("")
	assert.NotNil(t, app, "App should not be nil")
	assert.NoError(t, err, "Should not be err")
}

func TestNewAppWithInvalidConfigFile(t *testing.T) {
	dir := fs.NewDir(t, "test-config-file", fs.WithFile(".cinit.yaml", `xxx`))
	defer dir.Remove()
	app, err := NewApp(dir.Path() + "/.cinit.yaml")
	assert.Nil(t, app, "App should be nil")
	assert.Error(t, err, "Should be err")
	assert.Regexp(t, regexp.MustCompile("While parsing config: yaml: unmarshal errors"), err.Error(), "Should be same")
}