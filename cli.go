package main

import (
	"flag"
	"fmt"
	"os"
	"text/tabwriter"
	"time"
)

type command = int

const (
	COMMAND_ADD command = iota
	COMMAND_LIST
	COMMAND_SUMMARY
	COMMAND_DELETE
	COMMAND_EDIT
)

var commands = [...]string{"add", "list", "summary", "delete", "edit"}

func handleAdd(args []string) error {
	fs := flag.NewFlagSet("add", flag.ExitOnError)
	description := fs.String("description", "", "what did you pay for?")
	amount := fs.Float64("amount", 0, "how much did you pay?")

	if err := fs.Parse(args); err != nil {
		return fmt.Errorf("handleAdd: %w\n", err)
	}
	if *description == "" || *amount == 0 {
		fs.PrintDefaults()
		return nil
	}

	repo, err := newRepository()
	if err != nil {
		return fmt.Errorf("handleAdd: %w", err)
	}
	expense, err := newExpense(&repo, *description, *amount)
	if err != nil {
		return fmt.Errorf("handleAdd: %w", err)
	}
	repo.add(expense)
	repo.save()

	fmt.Printf("Expense added successfully (ID: %d)\n", expense.ID)
	return nil
}

func handleList() error {
	repo, err := newRepository()
	if err != nil {
		return fmt.Errorf("handleList: %w", err)
	}
	if len(repo.expenses) == 0 {
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 4, 0, '\t', 0)
	fmt.Fprintf(w, "ID\tDate\tDescription\tAmount\n")
	for _, v := range repo.expenses {
		fmt.Fprintf(w, "%d\t%s\t%s\t%.2f\n", v.ID, v.Date.Format(DATE_FORMAT), v.Description, v.Amount)
	}

	w.Flush()
	return nil
}

func handleSummary(args []string) error {
	var total float64

	fs := flag.NewFlagSet("summary", flag.ExitOnError)
	monthNumber := fs.Int("month", 0, "number of month - optional")

	if err := fs.Parse(args); err != nil {
		return fmt.Errorf("handleSummary: %w\n", invalidUsage)
	}

	if *monthNumber < 0 || *monthNumber > 12 {
		fs.PrintDefaults()
		return nil
	}

	repo, err := newRepository()
	if err != nil {
		return fmt.Errorf("handleSummary: %w", err)
	}

	month := time.Month(*monthNumber)
	fmt.Printf("Total expenses")

	if month != 0 {
		fmt.Printf(" for %s", month.String())
		for _, expense := range repo.expenses {
			if expense.Date.Month() != month {
				continue
			}
			total += expense.Amount
		}
	} else {
		for _, expense := range repo.expenses {
			total += expense.Amount
		}
	}

	fmt.Printf(": %.2f\n", total)
	return nil
}

func handleDelete(args []string) error {
	return nil
}

func handleEdit(args []string) error {
	return nil
}

func handleCommand(args []string) error {
	command := args[1]
	subArgs := args[2:]

	switch command {
	case commands[COMMAND_ADD]:
		return handleAdd(subArgs)
	case commands[COMMAND_LIST]:
		return handleList()
	case commands[COMMAND_SUMMARY]:
		return handleSummary(subArgs)
	case commands[COMMAND_DELETE]:
		return handleDelete(subArgs)
	case commands[COMMAND_EDIT]:
		return handleEdit(subArgs)
	default:
		return fmt.Errorf("handleCommand: %w: %s", invalidUsage, command)
	}
}
