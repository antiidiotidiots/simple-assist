package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/robertkrimen/otto"
)

type IntegrationConfiguration struct {
	Keywords []string `json:"keywords"`
	Script   string   `json:"script"`
}

var (
	WarningLogger *log.Logger
	InfoLogger    *log.Logger
	ErrorLogger   *log.Logger
)

var (
	currentDirectory, _   = os.Getwd()
	intergrationDirectory = filepath.Join(currentDirectory, "integrations")
)
var commandFlag = flag.String("command", "", "Command to run")

func init() {
	// LstdFlags is Ldate | Ltime
	// InfoLogger = log.New(os.Stdout, "INFO: ", log.Lshortfile)
	InfoLogger = log.New(os.Stdout, "INFO: ", 0)
	WarningLogger = log.New(os.Stderr, "WARNING: ", 0)
	ErrorLogger = log.New(os.Stderr, "ERROR: ", 0)
}

func main() {
	rand.Seed(time.Now().Unix())

	flag.Parse()
	// Get input from the user ( What are they telling the assistent? )
	var command string

	if *commandFlag == "" {
		fmt.Print("Ask me anything: ")
		command = singleLineInput()
	} else {
		command = *commandFlag
	}

	// Extract keywords
	extractedKeywords := extractKeywords(command)

	// fmt.Println(extractedKeywords)

	findKeywordsAndRun(extractedKeywords)
}

func extractKeywords(command string) []string {
	// Split the command by spaces
	splitCommand := strings.Split(command, " ")
	// splitCommand is of type []string

	// Remove all characters except a-z A-Z 0-9
	var charactersRegex = regexp.MustCompile(`[^a-zA-Z\d]`)

	// Loop through all the strings
	for index, word := range splitCommand {
		// Loop through all the characters in the string
		splitCommand[index] = charactersRegex.ReplaceAllString(word, "")
	}

	return splitCommand
}

func checkNilErr(err any) {
	if err != nil {
		// log.Fatalln("Error:\n%v\n", err)
		ErrorLogger.Fatalln(err)
	}
}

func singleLineInput() string {
	reader := bufio.NewReader(os.Stdin)

	input, err := reader.ReadString('\n')

	checkNilErr(err)

	input = strings.TrimSpace(input)

	return input
}

func loadConfiguration(path string) IntegrationConfiguration {
	var configuration IntegrationConfiguration
	configFile, err := os.OpenFile(path, os.O_RDONLY, 0644)
	checkNilErr(err)
	json.NewDecoder(configFile).Decode(&configuration)
	return configuration
}

func findKeywordsAndRun(extractedKeywords []string) {
	// Get all integrations inside the integrations folder
	folders, err := os.ReadDir(intergrationDirectory)
	checkNilErr(err)

	for _, f := range folders {
		if f.IsDir() {
			configFile := filepath.Join(intergrationDirectory, f.Name(), "config.json")

			configContents := loadConfiguration(configFile)

			integrationKeywords := configContents.Keywords
			for _, integrationKeyword := range integrationKeywords {
				// Check if integrationKeyword is any keyword we found earlier
				for _, extractedKeyword := range extractedKeywords {
					if integrationKeyword == extractedKeyword {
						scriptFile := filepath.Join(intergrationDirectory, f.Name(), configContents.Script)

						runScript(scriptFile, extractedKeywords)

						return
					}
				}
			}
		}
	}

	unknownCommand()
}

func runScript(scriptFile string, extractedKeywords []string) {
	scriptFileContents, err := os.ReadFile(scriptFile)
	checkNilErr(err)

	vm := otto.New()

	vm.Set("console.log", func(call otto.FunctionCall) otto.Value {
		InfoLogger.Printf("%s\n", call.Argument(0).String())
		return otto.Value{}
	})
	vm.Set("console.error", func(call otto.FunctionCall) otto.Value {
		ErrorLogger.Printf("%s\n", call.Argument(0).String())
		return otto.Value{}
	})
	vm.Set("console.warn", func(call otto.FunctionCall) otto.Value {
		WarningLogger.Printf("%s\n", call.Argument(0).String())
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

	checkNilErr(err)

	vm.Set("keywords", value)

	// fmt.Println(string(scriptFileContents))

	message, err := vm.Run(string(scriptFileContents))
	checkNilErr(err)

	fmt.Println(message)
}

func unknownCommand() {
	unknownCommandMessages := []string{
		"Sorry, I don't know how to answer that.",
		"I'm not sure.",
	}

	fmt.Println(unknownCommandMessages[rand.Intn(len(unknownCommandMessages))])
}
