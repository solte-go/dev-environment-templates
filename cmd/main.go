package main

import (
	"fmt"
	"kafka-scram-sasl/soltesandbox/internal/config"
	"kafka-scram-sasl/soltesandbox/internal/console"
	"kafka-scram-sasl/soltesandbox/internal/templates"
	"log"
	"os"
)

func main() {
	// load the template from a file
	conf, err := config.LoadConf()
	if err != nil {
		log.Fatal(err)
	}

	tmpl, err := templates.New(conf)
	if err != nil {
		log.Fatal(err)
	}

	console.Menu(tmpl)

	renderedServices := tmpl.PrepareChosenServices()

	outFile, err := os.Create("docker-compose.yaml")
	if err != nil {
		panic(err)
	}
	defer outFile.Close()

	fmt.Printf("%+v\n", renderedServices)

	//value, ok := data.(reflect.Value)
	//if !ok {
	//	value = reflect.ValueOf(data)
	//}

	//err = tmpl.Base.Execute(outFile, map[string]interface{}{"Services": renderedServices})
	err = tmpl.Base.Execute(outFile, map[string]interface{}{"Template": renderedServices})
	if err != nil {
		panic(err)
	}
}

//func parseTemplates(serviceNames []string, cfg *config.Config) map[string]string {
//	services := make(map[string]string)
//	for _, name := range serviceNames {
//		tmpl, err := template.ParseFiles(path(name + ".go.tmpl"))
//		if err != nil {
//			panic(err)
//		}
//		// Buffer to hold the template execution result
//		var tpl bytes.Buffer
//
//		// Execute the template and add indents
//		if err := tmpl.Execute(&tpl, cfg); err != nil {
//			panic(err)
//		}
//		// Add an indent of 2 spaces to each line
//		result := strings.ReplaceAll(tpl.String(), "\n", "\n  ")
//
//		services[name] = result
//	}
//
//	return services
//}
