package core

import (
	"regexp"
	"strings"
)

func ExtractKeywords(command string) []string {
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
