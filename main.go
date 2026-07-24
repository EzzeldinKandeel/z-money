package main

import (
	"errors"
	"fmt"
	"os"
)

const DATE_FORMAT string = "2006-01-02"

var DATA_FILE_PATH = os.Getenv("HOME") + "/.local/share/z-money/data"

var recordParsingError = errors.New("failed to parse csv record")
var indexOutOfBoundsError = errors.New("thid index does not exist in records array")
var recordDoesNotExist = errors.New("this record does not exist")

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
	arg := os.Args
	fmt.Println("arguments: ")
	fmt.Println(arg)
}
