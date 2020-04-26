package providers

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"text/template"

	log "github.com/sirupsen/logrus"
)

type DefaultTemplateProvider struct {
}

func NewDefaultTemplateProvider() *DefaultTemplateProvider {
	return &DefaultTemplateProvider{}
}

func (t *DefaultTemplateProvider) ResolveTemplates(source string, dest string, vars map[string]string) error {

	var srcContent []byte
	var err error

	if strings.HasPrefix(source, "env:") {
		// Build a template with just the env variable in it
		srcContent = []byte(fmt.Sprintf("{{ .%s }}", source[4:]))
	} else {
		// Create a new template and parse the letter into it.
		srcContent, err = ioutil.ReadFile(source)
		if err != nil {
			log.WithFields(log.Fields{
				"source": source,
				"dest":   dest,
			}).Error(err.Error())
			return err
		}
	}
	tmpl := template.Must(template.New("src").Parse(string(srcContent)))
	f, err := os.Create(dest)
	defer f.Close()
	if err != nil {
		log.Printf("Could not create template File %s", err.Error())
		return err
	}
	err = tmpl.Execute(f, vars)
	if err != nil {
		log.Printf("%s", err.Error())
		return err
	}

	return nil
}
