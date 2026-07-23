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
		return nil, err
	}

	dataReader := csv.NewReader(bytes.NewReader(data))
	dataReader.FieldsPerRecord = 4

	for {
		record, err := dataReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		expense, err := parseCsvRecord(record)
		if err != nil {
			return nil, err
		}
		expenses = append(expenses, expense)
	}
	return expenses, nil
}

func parseCsvRecord(record []string) (expense, error) {
	if len(record) != 4 {
		return expense{}, recordParsingError
	}

	var (
		rawId       = record[0]
		rawDate     = record[1]
		description = record[2]
		rawAmount   = record[3]
	)

	id, err := strconv.Atoi(rawId)
	if err != nil {
		return expense{}, recordParsingError
	}
	date, err := time.Parse(DATE_FORMAT, rawDate)
	if err != nil {
		return expense{}, recordParsingError
	}
	amount, err := strconv.ParseFloat(rawAmount, 64)
	if err != nil {
		return expense{}, recordParsingError
	}

	return expense{
		ID:          id,
		Description: description,
		Date:        date,
		Amount:      amount,
	}, nil
}

	fmt.Printf("data:\n%s", string(data))

}
