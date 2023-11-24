package tables

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
)

type Table struct {
	table table.Writer
}

// Creates a new table
func NewTable() Table {
	t := table.NewWriter()
	return Table{
		table: t,
	}
}

// Sets the header of the table
func (t *Table) SetHeader(header ...interface{}) {
	t.table.AppendHeader(table.Row(header))
}

// Adds a row to the table
func (t *Table) AddRow(row ...interface{}) {
	t.table.AppendRow(table.Row(row))
}

// Adds a separator between rows
func (t *Table) AddSeparator() {
	t.table.AppendSeparator()
}

// Enables auto-merging of cells with similar values in the given columns
func (t *Table) EnableAutoMergeOnColumns(columns ...int) {
	var colConfigs []table.ColumnConfig
	for _, c := range columns {
		colConfigs = append(colConfigs, table.ColumnConfig{Number: c, AutoMerge: true})
	}
	t.table.SetColumnConfigs(colConfigs)
}

// Renders the table
func (t *Table) Render() error {
	t.table.SetStyle(table.StyleLight)
	t.table.Style().Options.DrawBorder = false
	t.table.Style().Options.SeparateRows = false
	t.table.Style().Options.SeparateColumns = true
	t.table.Style().Options.SeparateHeader = true

	renderedTable := fmt.Sprintf("\n%s\n\n", t.table.Render())

	lessCmd := exec.Command("less", "-F", "-S", "-w")
	lessCmd.Stdin = strings.NewReader(renderedTable)
	lessCmd.Stdout = os.Stdout

	err := lessCmd.Run()
	if err != nil {
		return fmt.Errorf("run less command: %w", err)
	}

	return nil
}
