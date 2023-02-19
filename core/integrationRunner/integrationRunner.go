package core

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"

	"github.com/robertkrimen/otto"
	"github.com/sirupsen/logrus"

	request "github.com/antiidiotidiots/simple-assist/core"
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
					if strings.EqualFold(integrationKeyword, extractedKeyword) {
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

	// Setting some global functions that bind to Go functions
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

	//  Create assist object to allow for communication between the script and the Go code
	assist, err := vm.Object("({})")
	if err != nil {
		return err
	}

	assist.Set("respond", func(call otto.FunctionCall) otto.Value {
		for _, argument := range call.ArgumentList {
			fmt.Print(argument.String())
			fmt.Print(" ")
		}
		fmt.Println()
		return otto.Value{}
	})

	assist.Set("request", func(call otto.FunctionCall) otto.Value {
		url := call.Argument(0).String()
		method := call.Argument(1).String()

		bodyObject := call.Argument(2).Object()
		body := objectToMapDeep(bodyObject)

		headersObject := call.Argument(3).Object()
		headers := objectToMap(headersObject)

		response, err := request.Request(url, method, body, headers)
		if err != nil {
			return otto.Value{}
		}

		// Return a string with the response
		responseString, err := vm.ToValue(response)
		if err != nil {
			return otto.Value{}
		}

		return responseString
	})

	keywordsJavascriptString := "(["

	for _, extractedKeyword := range extractedKeywords {
		keywordsJavascriptString += "'"
		keywordsJavascriptString += extractedKeyword
		keywordsJavascriptString += "',"
	}

	keywordsJavascriptString += "])"

	// Example:
	// keywordsJavascriptString := "(['hello', 'world'])"

	// Get the extractedKeywords
	value, err := vm.Object(keywordsJavascriptString)
	if err != nil {
		return err
	}

	assist.Set("keywords", value)

	// Set the assist object
	vm.Set("assist", assist)

	_, err = vm.Run(string(scriptFileContents))
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func objectToMapDeep(object *otto.Object) map[string]interface{} {
	objectMap := make(map[string]interface{})

	objectKeys := object.Keys()
	for _, key := range objectKeys {
		value, _ := object.Get(key)
		if value.Class() == "Object" {
			objectMap[key] = objectToMapDeep(value.Object())
		} else {
			objectMap[key] = value.String()
		}
	}

	return objectMap
}

func objectToMap(object *otto.Object) map[string]string {
	objectMap := make(map[string]string)

	objectKeys := object.Keys()
	for _, key := range objectKeys {
		value, _ := object.Get(key)
		objectMap[key] = value.String()
	}

	return objectMap
}

func unknownCommand() {
	unknownCommandMessages := []string{
		"Sorry, I don't know how to answer that.",
		"I'm not sure.",
		"I'm sorry, I can't do that yet.",
	}

	fmt.Println(unknownCommandMessages[rand.Intn(len(unknownCommandMessages))])
}
