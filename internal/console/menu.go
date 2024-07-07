package console

import (
	"fmt"
	"kafka-scram-sasl/soltesandbox/internal/templates"
)

func Menu(tmpl *templates.Template) {
	// Menu loop

	var availableServices = make(map[int]string)

	for k, srv := range tmpl.Services {
		availableServices[k] = srv.Name
	}

	for {
		fmt.Println("\nAvailable Templates:")
		for i, name := range availableServices {
			fmt.Printf("[%d] %s\n", i, name)
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
			if value, ok := tmpl.Services[choice]; ok {
				//chosenTemplates = append(chosenTemplates, value)
				value.Build = true

				fmt.Printf("Template added: %v\n", value.Name)
				delete(availableServices, choice)
			}

		} else {
			fmt.Println("Invalid choice, please enter a number from the list")
		}

	}
}
