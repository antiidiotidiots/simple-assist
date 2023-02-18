package core

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"

	"github.com/robertkrimen/otto"
	"github.com/sirupsen/logrus"
)

type IntegrationConfiguration struct {
	Keywords []string `json:"keywords"`
	Script   string   `json:"script"`
}

var (
	currentDirectory, _   = os.Getwd()
	intergrationDirectory = filepath.Join(currentDirectory, "integrations")
)

func RunMatchingIntegration(extractedKeywords []string) error {
	// Get all integrations inside the integrations folder
	folders, err := os.ReadDir(intergrationDirectory)
	if err != nil {
		return err
	}

	for _, f := range folders {
		if f.IsDir() {
			configFile := filepath.Join(intergrationDirectory, f.Name(), "config.json")

			configContents, err := loadConfiguration(configFile)
			if err != nil {
				return err
			}

			integrationKeywords := configContents.Keywords
			for _, integrationKeyword := range integrationKeywords {
				// Check if integrationKeyword is any keyword we found earlier
				for _, extractedKeyword := range extractedKeywords {
					if integrationKeyword == extractedKeyword {
						scriptFile := filepath.Join(intergrationDirectory, f.Name(), configContents.Script)

						runScript(scriptFile, extractedKeywords)

						return nil
					}
				}
			}
		}
	}

	unknownCommand()
	return nil
}

func loadConfiguration(path string) (IntegrationConfiguration, error) {
	var configuration IntegrationConfiguration
	configFile, err := os.OpenFile(path, os.O_RDONLY, 0644)

	if err != nil {
		return configuration, err
	}

	json.NewDecoder(configFile).Decode(&configuration)
	return configuration, nil
}

func runScript(scriptFile string, extractedKeywords []string) error {
	scriptFileContents, err := os.ReadFile(scriptFile)
	if err != nil {
		return err
	}

	vm := otto.New()

	vm.Set("console.log", func(call otto.FunctionCall) otto.Value {
		logrus.Info("%s\n", call.Argument(0).String())
		return otto.Value{}
	})
	vm.Set("console.error", func(call otto.FunctionCall) otto.Value {
		logrus.Error("%s\n", call.Argument(0).String())
		return otto.Value{}
	})
	vm.Set("console.warn", func(call otto.FunctionCall) otto.Value {
		logrus.Warn("%s\n", call.Argument(0).String())
		return otto.Value{}
	})

	vm.Set("sayHello", func(call otto.FunctionCall) otto.Value {
		fmt.Printf("Hello, %s.\n", call.Argument(0).String())
		return otto.Value{}
	})

	keywordsJavascriptString := "(["

	for _, extractedKeyword := range extractedKeywords {
		keywordsJavascriptString += "'"
		keywordsJavascriptString += extractedKeyword
		keywordsJavascriptString += "',"
	}

	keywordsJavascriptString += "])"

	// Get the extractedKeywords
	value, err := vm.Object(keywordsJavascriptString)
	if err != nil {
		return err
	}

	vm.Set("keywords", value)

	// fmt.Println(string(scriptFileContents))

	message, err := vm.Run(string(scriptFileContents))
	if err != nil {
		return err
	}

	fmt.Println(message)
	return nil
}

func unknownCommand() {
	unknownCommandMessages := []string{
		"Sorry, I don't know how to answer that.",
		"I'm not sure.",
		"I'm sorry, I can't do that yet.",
	}

	fmt.Println(unknownCommandMessages[rand.Intn(len(unknownCommandMessages))])
}
