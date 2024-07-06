package console

import (
	"fmt"
	"os"
	"strings"
)

func Menu(dir string) []string {
	// Get all available templates
	files, err := os.ReadDir(dir)
	if err != nil {
		panic(err)
	}

	var templates = make(map[int]string)
	for i, file := range files {
		filename := file.Name()
		if strings.HasSuffix(filename, ".go.tmpl") && filename != "base.go.tmpl" {
			choice := strings.TrimSuffix(filename, ".go.tmpl")
			templates[i] = choice
		}
	}

	// Menu loop
	var chosenTemplates []string
	for {
		fmt.Println("\nAvailable Templates:")
		for i, tmpl := range templates {
			fmt.Printf("[%d] %s\n", i, tmpl)
		}
		fmt.Println("[0] Done")

		var choice int
		fmt.Print("\nEnter the number of the template you want to add (or 0 if you're done): ")
		_, err := fmt.Scan(&choice)
		if err != nil {
			fmt.Println("Invalid input, please enter a number")
			continue
		}

		if choice == 0 {
			break
		} else if choice > 0 {
			if value, ok := templates[choice]; ok {
				chosenTemplates = append(chosenTemplates, value)
				fmt.Printf("Template added: %v\n", value)
				delete(templates, choice)
			}

		} else {
			fmt.Println("Invalid choice, please enter a number from the list")
		}

	}
	return chosenTemplates
}
