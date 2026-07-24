package main

import (
	"errors"
	"fmt"
	"log"
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

type repository interface {
	load() error
	save() error
	add(record any)
	remove(index int) error
	update(index int, description string, amount float64) error
	find(id int) (int, error)
	getNextId() int
}

func main() {
	if err := os.MkdirAll(DATA_FILE_DIR, 0700); err != nil {
		log.Fatal(err)
	}
	arg := os.Args
	fmt.Println("arguments: ")
	fmt.Println(arg)
}
