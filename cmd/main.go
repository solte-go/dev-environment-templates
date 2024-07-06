package main

import (
	"bytes"
	"kafka-scram-sasl/soltesandbox/internal/config"
	"kafka-scram-sasl/soltesandbox/internal/console"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

const tmplDir = "./internal/templates"

func main() {
	// load the template from a file
	cfg, err := config.LoadConf()
	if err != nil {
		log.Fatal(err)
	}

	chosenTmpl := console.Menu(tmplDir)

	renderedServices := parseTemplates(chosenTmpl, cfg)

	baseTmpl, err := template.ParseFiles(path("base.go.tmpl"))
	if err != nil {
		panic(err)
	}

	outFile, err := os.Create("docker-compose.yaml")
	if err != nil {
		panic(err)
	}
	defer outFile.Close()

	err = baseTmpl.Execute(outFile, map[string]interface{}{"Services": renderedServices})
	if err != nil {
		panic(err)
	}
}

func path(filename string) string {
	return filepath.Join(tmplDir, filename)
}

func parseTemplates(serviceNames []string, cfg *config.Config) map[string]string {
	services := make(map[string]string)
	for _, name := range serviceNames {
		tmpl, err := template.ParseFiles(path(name + ".go.tmpl"))
		if err != nil {
			panic(err)
		}
		// Buffer to hold the template execution result
		var tpl bytes.Buffer

		// Execute the template and add indents
		if err := tmpl.Execute(&tpl, cfg); err != nil {
			panic(err)
		}
		// Add an indent of 2 spaces to each line
		result := strings.ReplaceAll(tpl.String(), "\n", "\n  ")

		services[name] = result
	}

	return services
}
