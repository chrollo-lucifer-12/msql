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
		os.Exit(1)
	}

	database, err := db.FromFile(os.Args[1])
	if err != nil {
		os.Exit(1)
	}

	err = cli(database)
	if err != nil {
		os.Exit(1)
	}
}

func cli(database *db.Db) error {
	if err := printFlushed("mqlite> "); err != nil {
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

		if err := printFlushed("\nmqlite> "); err != nil {
			return err
		}
	}

	return nil
}

func displayTables(database *db.Db) error {

	scanner := database.GetScanner(1)

	for {

		cursor, err := scanner.NextRecord()

		if err != nil {
			break
		}

		if cursor == nil {
			break
		}

		typeValue := cursor.Field(0)

		if typeValue == nil {
			continue
		}

		value, ok := typeValue.AsString()

		if !ok {
			continue
		}

		if value == "table" {

			nameValue := cursor.Field(1)

			if nameValue == nil {
				continue
			}

			name, ok := nameValue.AsString()

			if ok {
				fmt.Printf("%s ", name)
			}
		}
	}

	return nil
}

func printFlushed(s string) error {
	_, err := fmt.Print(s)
	if err != nil {
		return fmt.Errorf("write stdout: %w", err)
	}

	return nil
}
