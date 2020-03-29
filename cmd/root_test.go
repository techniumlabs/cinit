package cmd

import (
	"path/filepath"
	"testing"

	"os"

	"io/ioutil"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"gotest.tools/fs"
)

func TestConfigFileInSameDirectory(t *testing.T) {
	dir := fs.NewDir(t, "test-config-file", fs.WithFile(".cinit.yaml", ""))
	defer dir.Remove()
	os.Chdir(dir.Path())
	initConfig()
	fexp, _ := filepath.EvalSymlinks(dir.Path() + "/.cinit.yaml")
	factual, _ := filepath.EvalSymlinks(viper.GetViper().ConfigFileUsed())
	assert.Equal(t, fexp, factual, "File Path should be equal")
}

func TestConfigFileInHomeDirectory(t *testing.T) {
	home, _ := os.UserHomeDir()
	ioutil.WriteFile(home+"/.cinit.yaml", nil, 0644)
	initConfig()
	assert.Equal(t, home+"/.cinit.yaml", viper.GetViper().ConfigFileUsed(), "File Path should be equal")
}

func TestConfigFileUserProvided(t *testing.T) {
	dir := fs.NewDir(t, "test-config-file", fs.WithFile(".cinit.yam", ""))
	defer dir.Remove()
	cfgFile = dir.Path() + "/.cinit.yaml"
	initConfig()
	fexp, _ := filepath.EvalSymlinks(dir.Path() + "/.cinit.yaml")
	factual, _ := filepath.EvalSymlinks(viper.GetViper().ConfigFileUsed())
	assert.Equal(t, fexp, factual, "File Path should be equal")
}
