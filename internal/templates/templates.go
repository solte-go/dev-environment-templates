package templates

import (
	"bytes"
	"fmt"
	"kafka-scram-sasl/soltesandbox/internal/config"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

const (
	baseTmpl    = "./internal/templates/base/base.go.tmpl"
	servicesDir = "./internal/templates/service"
	volumeDir   = "./internal/templates/volume"
	tmplExt     = ".go.tmpl"
)

type Service struct {
	Name          string
	Configuration string
	Volume        string
	Network       string
	Build         bool
}

type PrepareData struct {
	Service []string
	Volume  []string
}

type Template struct {
	Base     *template.Template
	Services map[int]*Service
}

func New(conf *config.Config) (*Template, error) {
	var tmpl = &Template{}

	if err := tmpl.base(); err != nil {
		return nil, err
	}

	if err := tmpl.buildConfigs(conf); err != nil {
		return nil, err
	}

	return tmpl, nil
}

func (t *Template) base() error {
	var err error
	t.Base, err = template.ParseFiles(baseTmpl)
	if err != nil {
		return fmt.Errorf("base template parsing failed: %w", err)
	}
	return nil
}

func (t *Template) PrepareChosenServices() PrepareData {
	var result PrepareData

	for _, s := range t.Services {
		if s.Build {
			result.Service = append(result.Service, s.Configuration)
			if s.Volume != "" {
				result.Volume = append(result.Volume, s.Volume)
			}
		}
	}

	return result
}

func (t *Template) buildConfigs(conf *config.Config) error {

	t.Services = make(map[int]*Service)

	// Get all available templates
	files, err := os.ReadDir(servicesDir)
	if err != nil {
		return fmt.Errorf("failed to load confguration templates: %w", err)
	}

	for i, file := range files {
		filename := file.Name()
		if strings.HasSuffix(filename, tmplExt) && filename != "base.go.tmpl" {
			choice := strings.TrimSuffix(filename, tmplExt)

			srvConfig, err := parseTemplates(servicesDir, choice, conf)
			if err != nil {
				return err
			}

			t.Services[i+1] = &Service{
				Name:          choice,
				Configuration: srvConfig,
			}
		}
	}

	files, err = os.ReadDir(volumeDir)
	if err != nil {
		return fmt.Errorf("failed to load volume templates: %w", err)
	}

	for _, file := range files {
		filename := file.Name()

		for _, srv := range t.Services {
			fmt.Println(srv.Name)
		}

		if strings.HasSuffix(filename, tmplExt) && filename != "base.go.tmpl" {
			filename = strings.TrimSuffix(filename, tmplExt)
			for _, s := range t.Services {
				if s.Name == filename {
					s.Volume, err = parseTemplates(volumeDir, filename, conf)
					if err != nil {
						return fmt.Errorf("failed to load volume templates: %w", err)
					}
				}
			}
		}
	}

	return nil
}

func parseTemplates(dir, name string, cfg *config.Config) (string, error) {
	//services := make(map[string]string)

	tmpl, err := template.ParseFiles(path(dir, name))
	if err != nil {
		return "", fmt.Errorf("failed to parse templates: %w", err)
	}
	// Buffer to hold the template execution result
	var tpl bytes.Buffer

	// Execute the template and add indents
	if err := tmpl.Execute(&tpl, cfg); err != nil {
		return "", fmt.Errorf("failed to execute templates: %w", err)
	}
	// Add an indent of 2 spaces to each line
	//result := strings.ReplaceAll(tpl.String(), "\n", "\n  ")

	return strings.ReplaceAll(tpl.String(), "\n", "\n  "), nil
}

func path(dir, filename string) string {
	return filepath.Join(dir, filename+tmplExt)
}
