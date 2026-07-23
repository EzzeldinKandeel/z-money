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
	t := time.Now()
	fmt.Printf("date: %s\n", t.Format("2006-01-02"))
}

func createExpense(description string, amount float64) (expense, error) {
	if amount <= 0 {
		return expense{}, errors.New("amount must be a positive number")
	}

	savedExpenses, err := loadExpenses()
	if err != nil {
		return expense{}, err
	}

	return expense{
		ID:          getNextId(savedExpenses),
		Date:        time.Now(),
		Description: description,
		Amount:      amount,
	}, nil
}

func loadExpenses() ([]expense, error) {
	var expenses []expense

	data, err := os.ReadFile(DATA_FILE_PATH)
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
