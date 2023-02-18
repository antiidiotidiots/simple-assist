package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	integrationRunner "github.com/antiidiotidiots/simple-assist/core/integrationRunner"
	keywordExtractor "github.com/antiidiotidiots/simple-assist/core/keywordExtractor"

	"github.com/alexflint/go-arg"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

var args struct {
	Command  []string `arg:"positional" help:"Command to run"`
	Repeat   bool     `arg:"-r, --repeat, env:REPEAT" help:"Repeat the command after it is done running" default:"false"`
	LogLevel string   `arg:"--log-level, env:LOG_LEVEL" help:"\"debug\", \"info\", \"warning\", \"error\", or \"fatal\"" default:"info"`
	LogColor bool     `arg:"--log-color, env:LOG_COLOR" help:"Force colored logs" default:"false"`
}

func main() {
	godotenv.Load(".env")
	arg.MustParse(&args)

	logrus.SetOutput(os.Stdout)
	logrus.SetFormatter(&logrus.TextFormatter{PadLevelText: true, DisableQuote: true, ForceColors: args.LogColor, DisableColors: !args.LogColor})
	if args.LogLevel == "debug" {
		logrus.SetLevel(logrus.DebugLevel)
		// Enable line numbers in debug logs - Doesn't help too much since a fatal error still needs to be debugged
		logrus.SetReportCaller(true)
	} else if args.LogLevel == "info" {
		logrus.SetLevel(logrus.InfoLevel)
	} else if args.LogLevel == "warning" {
		logrus.SetLevel(logrus.WarnLevel)
	} else if args.LogLevel == "error" {
		logrus.SetLevel(logrus.ErrorLevel)
	} else if args.LogLevel == "fatal" {
		logrus.SetLevel(logrus.FatalLevel)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}

	// Get input from the user ( What are they telling the assistant? )
	command := strings.Join(args.Command, " ")
	var err error
	var repeat bool
	if command == "" {
		err = askQuestion()
		repeat = true
	} else {
		err = run(command)
	}

	if err != nil {
		logrus.Fatal(err)
	}

	for args.Repeat || repeat {
		err := askQuestion()
		if err != nil {
			logrus.Fatal(err)
		}
	}
}

func askQuestion() error {
	fmt.Print("Ask me anything: ")
	command, err := singleLineInput()
	if err != nil {
		logrus.Fatal(err)
	}

	err = run(command)
	if err != nil {
		return err
	}

	return nil
}

func run(command string) error {
	// Extract keywords
	extractedKeywords := keywordExtractor.ExtractKeywords(command)

	err := integrationRunner.RunMatchingIntegration(extractedKeywords)
	if err != nil {
		return err
	}

	return nil
}

func singleLineInput() (string, error) {
	reader := bufio.NewReader(os.Stdin)

	input, err := reader.ReadString('\n')

	if err != nil {
		return "", err
	}

	input = strings.TrimSpace(input)

	return input, nil
}
