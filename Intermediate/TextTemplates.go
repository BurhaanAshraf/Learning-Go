package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"text/template"
)

func Templates() {

	// ============================================================
	// 🔹 Basic Template
	// ============================================================

	baseTemplate := template.New("Example")

	// Parse converts raw string → compiled template (*template.Template)
	baseTemplate, err := baseTemplate.Parse(
		"Welcome, {{.name}}! How are you\n",
	)

	if err != nil {
		panic(err)
	}

	baseData := map[string]interface{}{
		"name": "Burhaan",
	}

	// Execute replaces placeholder with actual value
	err = baseTemplate.Execute(os.Stdout, baseData)

	if err != nil {
		panic(err)
	}


	// ============================================================
	// 🔹 template.Must Shortcut
	// ============================================================

	// Must → panic automatically if parsing fails
	secondTemplate := template.Must(
		template.New("Example2").Parse(
			"Welcome again, {{.prevname}}! Good to see you back\n",
		),
	)

	secondData := map[string]interface{}{
		"prevname": "Burhaan",
	}

	secondTemplate.Execute(os.Stdout, secondData)


	// ============================================================
	// 🔹 User Input
	// ============================================================

	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Enter your name:")

	userName, _ := reader.ReadString('\n')
	userName = strings.TrimSpace(userName)


	// ============================================================
	// 🔹 Template Map (Raw Strings)
	// ============================================================

	templateMap := map[string]string{
		"welcome":      "Welcome, {{.name}}! We are glad you joined.",
		"notification": "{{.name}}, you have a new notification. {{.notification}}",
		"error":        "OOPS! Error occured {{.errormessage}}",
	}


	// ============================================================
	// 🔹 Parse Templates → Store as *template.Template
	// ============================================================

	// Here we convert raw template strings → compiled templates

	parsedMap := make(map[string]*template.Template)

	for tmplName, tmplStr := range templateMap {

		// After Parse → each value becomes *template.Template
		parsedMap[tmplName] =
			template.Must(template.New(tmplName).Parse(tmplStr))
	}


	// ============================================================
	// 🔹 Menu Loop
	// ============================================================

	for {

		fmt.Println("\nMenu")
		fmt.Println("1. Join")
		fmt.Println("2. Get notification")
		fmt.Println("3. Get Error")
		fmt.Println("4. Exit")
		fmt.Println("Choose an option")

		userChoice, _ := reader.ReadString('\n')
		userChoice = strings.TrimSpace(userChoice)

		var templateData map[string]interface{}
		var selectedTemplate *template.Template


		switch userChoice {

		// --------------------------------------------------------
		// Case 1: Welcome
		// → Uses compiled template to inject name
		// --------------------------------------------------------

		case "1":
			selectedTemplate = parsedMap["welcome"]

			templateData = map[string]interface{}{
				"name": userName,
			}


		// --------------------------------------------------------
		// Case 2: Notification
		// → Takes input and injects into template
		// --------------------------------------------------------

		case "2":

			fmt.Println("Enter your notification message")

			notificationMsg, _ := reader.ReadString('\n')
			notificationMsg = strings.TrimSpace(notificationMsg)

			selectedTemplate = parsedMap["notification"]

			templateData = map[string]interface{}{
				"name":         userName,
				"notification": notificationMsg,
			}


		// --------------------------------------------------------
		// Case 3: Error
		// → Takes error input and injects into template
		// --------------------------------------------------------

		case "3":

			fmt.Println("Enter Error Message")

			errorMsg, _ := reader.ReadString('\n')
			errorMsg = strings.TrimSpace(errorMsg)

			selectedTemplate = parsedMap["error"]

			templateData = map[string]interface{}{
				"errormessage": errorMsg,
			}


		// --------------------------------------------------------
		// Case 4: Exit
		// --------------------------------------------------------

		case "4":
			fmt.Println("Exiting...")
			return


		default:
			fmt.Println("Please choose a valid option!")
			continue
		}


		// ========================================================
		// 🔹 Execute Template
		// ========================================================

		// Runs compiled template with data
		err = selectedTemplate.Execute(os.Stdout, templateData)

		if err != nil {
			panic("Error executing template")
		}
	}
}


// ============================================================
// 🔹 Quick Revision
// ============================================================

// 🔹 What is *template.Template?

// *template.Template is a parsed (compiled) template object
// It understands placeholders like {{.name}}
// and can be executed with data

// Flow:
// Raw string → Parse → *template.Template → Execute → Output

/// *template.Template → compiled template object

// template.New().Parse() → creates *template.Template


// Execute() → replaces placeholders with data

// {{.key}} → dynamic placeholder

// Must() → panic if parsing fails

// Raw string → Parsed → Executed → Output

// Best Practice:
// - Parse templates once
// - Reuse them
// - Keep execution separate