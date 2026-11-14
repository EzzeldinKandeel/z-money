package main

import (
	"errors"
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
	if err = repo.save(); err != nil {
		return fmt.Errorf("handleAdd: %w", err)
	}

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
		total, _ = repo.getTotalForMonth(month)
	} else {
		total = repo.getTotal()
	}

	fmt.Printf(": %.2f\n", total)
	return nil
}

func handleDelete(args []string) error {
	fs := flag.NewFlagSet("delete", flag.ExitOnError)
	id := fs.Int("id", 0, "id of expense you want to delete")
	if err := fs.Parse(args); err != nil {
		return fmt.Errorf("handleDelete: %w\n", invalidUsage)
	}

	if *id < 1 {
		fs.PrintDefaults()
		return nil
	}

	repo, err := newRepository()
	if err != nil {
		return fmt.Errorf("handleDelete: %w", err)
	}
	index, err := repo.find(*id)
	if errors.Is(err, recordDoesNotExist) {
		fmt.Fprintf(os.Stderr, "%s", recordDoesNotExist.Error())
		os.Exit(1)
	}
	if err != nil {
		return fmt.Errorf("handleDelete: %w", err)
	}
	if err = repo.remove(index); err != nil {
		return fmt.Errorf("handleDelete: %w", err)
	}
	if err = repo.save(); err != nil {
		return fmt.Errorf("handleDelete: %w", err)
	}

	fmt.Println("Expense deleted successfully")
	return nil
}

func handleEdit(args []string) error {
	fs := flag.NewFlagSet("edit", flag.ExitOnError)
	id := fs.Int("id", 0, "id of the expense you want to edit")
	description := fs.String("description", "", "what did you pay for? - optional")
	amount := fs.Float64("amount", 0, "how much did you pay? - optional")

	if err := fs.Parse(args); err != nil {
		return fmt.Errorf("handleEdit: %w", err)
	}
	if *id < 1 {
		fs.PrintDefaults()
		return nil
	}

	repo, err := newRepository()
	if err != nil {
		return fmt.Errorf("handleEdit: %w", err)
	}
	index, err := repo.find(*id)
	if errors.Is(err, recordDoesNotExist) {
		fmt.Fprintf(os.Stderr, "%s", recordDoesNotExist.Error())
		os.Exit(1)
	}
	if err != nil {
		return fmt.Errorf("handleEdit: %w", err)
	}
	if err = repo.update(index, *description, *amount); err != nil {
		return fmt.Errorf("handleEdit: %w", err)
	}
	if err = repo.save(); err != nil {
		return fmt.Errorf("handleEdit: %w", err)
	}

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
