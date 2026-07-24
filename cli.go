package main

import (
	"os"
	"slices"
)

var ALLOWED_COMMANDS = [...]string{"add", "list", "summary", "delete", "edit"}

func parseCommand([]string) (string, error) {
	if len(os.Args) < 2 {
		return "", invalidUsage
	}

	command := os.Args[1]
	if !slices.Contains(ALLOWED_COMMANDS[:], command) {
		return "", invalidUsage
	}

	return command, nil
}
