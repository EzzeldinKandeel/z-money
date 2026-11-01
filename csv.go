package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"slices"
	"strconv"
	"time"
)

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

func overwriteExpenses(expenses [][4]string) error {
	dataFile, err := os.OpenFile(DATA_FILE_PATH, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("overwriteExpenses: %w", err)
	}
	defer dataFile.Close()

	var expenseSlices [][]string
	for e := range slices.Values(expenses) {
		expenseSlices = append(expenseSlices, e[:])
	}

	err = dataFile.Truncate(0)
	if err != nil {
		return fmt.Errorf("overwriteExpenses: %w", err)
	}
	dataWriter := csv.NewWriter(dataFile)
	if err = dataWriter.WriteAll(expenseSlices); err != nil {
		return fmt.Errorf("overwriteExpenses: %w", err)
	}

	return nil
}
