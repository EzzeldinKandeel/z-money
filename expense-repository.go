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
	"time"
)

type expensesRepository struct {
	expenses []expense
}

func newRepository() (expensesRepository, error) {
	repo := expensesRepository{}
	if err := repo.load(); err != nil {
		return expensesRepository{}, fmt.Errorf("newRepository: %w", err)
	}

	return repo, nil
}

func (r *expensesRepository) load() error {
	var expenses []expense

	data, err := os.ReadFile(DATA_FILE_PATH)
	if errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("expensesRespository.load: IT DOESNOT EXIST %w", err)
	}
	if err != nil {
		return fmt.Errorf("expensesRespository.load: %w", err)
	}

	dataReader := csv.NewReader(bytes.NewReader(data))
	dataReader.FieldsPerRecord = 4

	for {
		record, err := dataReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("expensesRespository.load: %w", err)
		}

		expense, err := parseCsvRecord(record)
		if err != nil {
			return fmt.Errorf("expensesRespository.load: %w", err)
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
		return fmt.Errorf("expensesRespository.save: %w", err)
	}

	return nil
}

func (r *expensesRepository) add(record expense) {
	r.expenses = append(r.expenses, record)
}

func (r *expensesRepository) remove(index int) error {
	if index >= len(r.expenses) {
		return fmt.Errorf("expensesRepository.remove: %w", indexOutOfBoundsError)
	}

	r.expenses = slices.Delete(r.expenses, index, index+1)
	return nil
}

func (r *expensesRepository) update(index int, description string, amount float64) error {
	if description != "" {
		r.expenses[index].Description = description
	}
	if amount < 0 {
		return fmt.Errorf("expensesRepository.update: %w", recordParsingError)
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
		return -1, fmt.Errorf("expensesRepository.find: %w", recordDoesNotExist)
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
		return 0, fmt.Errorf("expensesRepository.getTotalForMonth: month must be between 1 and 12")
	}

	total := 0.0
	for e := range slices.Values(r.expenses) {
		if e.Date.Month() == month {
			total += e.Amount
		}
	}

	return total, nil
}
