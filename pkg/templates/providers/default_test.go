package providers

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"gotest.tools/fs"
)

func TestNewDefaultTemplateProvider(t *testing.T) {
	provider := NewDefaultTemplateProvider()
	assert.NotNil(t, provider, "Default Provider should not be nil")
	assert.IsType(t, &DefaultTemplateProvider{}, provider, "Should be of Template Provider type")
}

func TestDefaultTemplateResolve(t *testing.T) {
	tpl := NewDefaultTemplateProvider()

	dir := fs.NewDir(t, "test-template", fs.WithFile("src-template", `{{ .TESTVAR1 }}`))
	defer dir.Remove()
	myvarmap := map[string]string{
		"TESTVAR1": "SOMEVALUE1",
	}
	err := tpl.ResolveTemplates(dir.Path()+"/src-template", dir.Path()+"/dest-file", myvarmap)
	assert.NoError(t, err, "Should not throw error")
	assert.FileExists(t, dir.Path()+"/dest-file", "Should have a file")
	destContent, err := ioutil.ReadFile(dir.Path() + "/dest-file")
	assert.NoError(t, err, "Should not have any error Reading file")
	assert.Equal(t, "SOMEVALUE1", string(destContent), "Content Should be same")
}

func TestDefaultTemplateResolveInvalidInput(t *testing.T) {
	tpl := NewDefaultTemplateProvider()

	dir := fs.NewDir(t, "test-template", fs.WithFile("src-template", `{{ .TESTVAR2 }}`))
	defer dir.Remove()
	myvarmap := map[string]string{
		"TESTVAR1": "SOMEVALUE1",
	}
	err := tpl.ResolveTemplates(dir.Path()+"/src-template", dir.Path()+"/dest-file", myvarmap)
	assert.NoError(t, err, "Should not throw error")
	assert.FileExists(t, dir.Path()+"/dest-file", "Should have a file")
	destContent, err := ioutil.ReadFile(dir.Path() + "/dest-file")
	assert.NoError(t, err, "Should not have any error Reading file")
	assert.Equal(t, "<no value>", string(destContent), "Content Should be same")
}

func TestDefaultTemplateResolveInvalidSrcFile(t *testing.T) {
	tpl := NewDefaultTemplateProvider()

	dir := fs.NewDir(t, "test-template")
	defer dir.Remove()
	myvarmap := map[string]string{
		"TESTVAR1": "SOMEVALUE1",
	}
	err := tpl.ResolveTemplates(dir.Path()+"/src-template", dir.Path()+"/dest-file", myvarmap)
	assert.EqualError(t, err, fmt.Sprintf("open %s/src-template: no such file or directory", dir.Path()), "Should throw error")
}
