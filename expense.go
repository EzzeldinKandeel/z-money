package main

import (
	"strconv"
	"time"
)

type expense struct {
	ID          int
	Date        time.Time
	Description string
	Amount      float64
}

func newExpense(repo repository[expense], description string, amount float64) (expense, error) {
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

func formatExpenseForSaving(expense expense) [4]string {
	var values [4]string
	values[0] = strconv.Itoa(expense.ID)
	values[1] = expense.Date.Format(DATE_FORMAT)
	values[2] = expense.Description
	values[3] = strconv.FormatFloat(expense.Amount, 'f', 2, 64)
	return values
}
