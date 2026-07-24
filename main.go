package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"slices"
	"strconv"
	"time"
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

type expense struct {
	ID          int
	Date        time.Time
	Description string
	Amount      float64
}

func main() {
	arg := os.Args
	fmt.Println("arguments: ")
	fmt.Println(arg)
}

func newExpense(repo repository, description string, amount float64) (expense, error) {
	if amount <= 0 {
		return expense{}, recordParsingError
	}

	return expense{
		ID:          repo.getNextId(),
		Date:        time.Now(),
		Description: description,
		Amount:      amount,
	}, nil
}

func newRepository() (expensesRepository, error) {
	repo := expensesRepository{}
	if err := repo.load(); err != nil {
		return expensesRepository{}, err
	}

	return repo, nil
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

func formatExpenseForSaving(expense expense) [4]string {
	var values [4]string
	values[0] = strconv.Itoa(expense.ID)
	values[1] = expense.Date.Format(DATE_FORMAT)
	values[2] = expense.Description
	values[3] = strconv.FormatFloat(expense.Amount, 'f', 2, 64)
	return values
}

func overwriteExpenses(expenses [][4]string) error {
	dataFile, err := os.OpenFile(DATA_FILE_PATH, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer dataFile.Close()

	var expenseSlices [][]string
	for e := range slices.Values(expenses) {
		expenseSlices = append(expenseSlices, e[:])
	}

	dataWriter := csv.NewWriter(dataFile)
	if err = dataWriter.WriteAll(expenseSlices); err != nil {
		return err
	}

	return nil
}
