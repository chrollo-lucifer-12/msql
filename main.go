package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/msql/db"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("missing db file")
		os.Exit(1)
	}

	database, err := db.FromFile(os.Args[1])
	if err != nil {
		fmt.Errorf("open database: %w", err)
		os.Exit(1)
	}

	err = cli(database)
	if err != nil {
		fmt.Errorf("cli error: %w", err)
		os.Exit(1)
	}
}

func cli(database *db.Db) error {
	if err := printFlushed("rqlite> "); err != nil {
		return err
	}

	reader := bufio.NewReader(os.Stdin)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}

		switch strings.TrimSpace(line) {
		case ".exit":
			return nil

		case ".tables":
			if err := displayTables(database); err != nil {
				return err
			}

		default:
			fmt.Printf("Unrecognized command '%s'\n", strings.TrimSpace(line))
		}

		if err := printFlushed("\nrqlite> "); err != nil {
			return err
		}
	}

	return nil
}

func displayTables(database *db.Db) error {
	scanner := database.GetScanner(1)

	for {
		record, err := scanner.NextRecord()

		if err != nil {
			fmt.Printf("error displaying: %w", err)
			break
		}

		typeValue := record.Field(0)

		value, ok := typeValue.AsString()
		if !ok {
			continue
		}

		if value == "table" {
			nameValue := record.Field(1)

			name, ok := nameValue.AsString()
			if ok {
				fmt.Printf("%s ", name)
			}
		}
	}

	return nil
}

func printFlushed(s string) error {
	fmt.Print(s)

	err := os.Stdout.Sync()
	if err != nil {
		return fmt.Errorf("flush stdout: %w", err)
	}

	return nil
}
