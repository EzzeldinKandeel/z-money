package main

import (
	"bytes"
	"cmp"
	"encoding/csv"
	"errors"
	"io"
	"os"
	"slices"
	"time"
)

type expensesRepository struct {
	expenses []expense
}

func (r *expensesRepository) load() error {
	var expenses []expense

	data, err := os.ReadFile(DATA_FILE_PATH)
	if err != nil {
		return err
	}

	dataReader := csv.NewReader(bytes.NewReader(data))
	dataReader.FieldsPerRecord = 4

	for {
		record, err := dataReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		expense, err := parseCsvRecord(record)
		if err != nil {
			return err
		}
		expenses = append(expenses, expense)
	}

	r.expenses = expenses
	return nil
}

func (r *expensesRepository) getNextId() int {
	if len(r.expenses) == 0 {
		return 1
	}
	lastExpense := slices.MaxFunc(r.expenses, func(a, b expense) int {
		return cmp.Compare(a.ID, b.ID)
	})

	return lastExpense.ID + 1
}

func (r *expensesRepository) save() error {
	var formatedExpenses [][4]string
	for e := range slices.Values(r.expenses) {
		formatedExpenses = append(formatedExpenses, formatExpenseForSaving(e))
	}

	if err := overwriteExpenses(formatedExpenses); err != nil {
		return err
	}

	return nil
}

func (r *expensesRepository) add(expense expense) {
	r.expenses = append(r.expenses, expense)
}

func (r *expensesRepository) remove(index int) error {
	if index >= len(r.expenses) {
		return indexOutOfBoundsError
	}

	r.expenses = slices.Delete(r.expenses, index, index+1)
	return nil
}

func (r *expensesRepository) update(index int, description string, amount float64) error {
	if description != "" {
		r.expenses[index].Description = description
	}
	if amount < 0 {
		return recordParsingError
	}
	if amount != 0 {
		r.expenses[index].Amount = amount
	}

	return nil
}

func (r *expensesRepository) find(id int) (int, error) {
	index := slices.IndexFunc(r.expenses, func(expense expense) bool {
		return expense.ID == id
	})
	if index == -1 {
		return -1, recordDoesNotExist
	}
	return index, nil
}

func (r *expensesRepository) getTotal() float64 {
	total := 0.0
	for e := range slices.Values(r.expenses) {
		total += e.Amount
	}

	return total
}

func (r *expensesRepository) getTotalForMonth(month time.Month) (float64, error) {
	if month < 1 || month > 12 {
		return 0, errors.New("month must be between 1 and 12")
	}

	total := 0.0
	for e := range slices.Values(r.expenses) {
		if e.Date.Month() == month {
			total += e.Amount
		}
	}

	return total, nil
}
