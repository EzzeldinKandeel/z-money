package main

import (
	"errors"
	"fmt"
	"os"
)

const DATE_FORMAT string = "2006-01-02"
const DATA_FILE_NAME string = "data"

var DATA_FILE_DIR string = os.Getenv("HOME") + "/.local/share/z-money/"
var DATA_FILE_PATH string = DATA_FILE_DIR + DATA_FILE_NAME

var recordParsingError = errors.New("failed to parse csv record")
var indexOutOfBoundsError = errors.New("thid index does not exist in records array")
var recordDoesNotExist = errors.New("this record does not exist")
var invalidUsage = errors.New("unsupported command or flag")

type repository[T any] interface {
	load() error
	save() error
	add(record T)
	remove(index int) error
	update(index int, description string, amount float64) error
	find(id int) (int, error)
	getNextId() int
}

func main() {
	if err := os.MkdirAll(DATA_FILE_DIR, 0700); err != nil {
		fmt.Fprintln(os.Stderr, fmt.Errorf("main: %w", err))
		os.Exit(1)
	}

	file, err := os.OpenFile(DATA_FILE_PATH, os.O_CREATE, 0700)
	if err != nil {
		fmt.Fprintln(os.Stderr, fmt.Errorf("main: %w", err))
		os.Exit(1)
	}
	file.Close()

	args := os.Args
	if len(args) < 2 {
		fmt.Fprintln(
			os.Stderr,
			fmt.Errorf(
				"%w\navilable commands: %v",
				fmt.Errorf("main: %w", invalidUsage),
				commands,
			),
		)
		os.Exit(1)
	}
	if err := handleCommand(args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
