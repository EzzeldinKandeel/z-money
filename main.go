package main

import (
	"bytes"
	"cmp"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"slices"
	"strconv"
	"time"
)

const DATE_FORMAT string = "2006-01-02"

var DATA_FILE_PATH = os.Getenv("HOME") + "/.local/share/z-money/data"

var recordParsingError = errors.New("failed to parse csv record")

type expense struct {
	ID          int
	Date        time.Time
	Description string
	Amount      float64
}

func main() {
	dataDir := os.Getenv("HOME")+"/.local/share/z-money"
	_, err := os.ReadDir(dataDir)
	if err != nil {
		log.Printf("! Data directory does not exist.\n%v", err)
		log.Printf("Creating data directory...\n")

		mkDirErr := os.Mkdir(dataDir, 0700)
		if mkDirErr != nil {
			log.Fatal("! Could not create data directory. Exiting...\n")
		}

		err = nil
		log.Printf("Data directory created successfully.")
	}

	dataFile, err := os.OpenFile(dataDir+"/data", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal("! Could not access data file. Exiting...\nError:\n%w", err)
	}
	data := make([]byte, 100)
	dataFile.Read(data)

	fmt.Printf("data:\n%s", string(data))

}
