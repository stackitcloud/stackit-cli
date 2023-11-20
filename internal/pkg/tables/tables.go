package tables

import (
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
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
func (t *Table) Render(cmd *cobra.Command) {
	t.table.SetStyle(table.StyleLight)
	t.table.Style().Options.DrawBorder = false
	t.table.Style().Options.SeparateRows = false
	t.table.Style().Options.SeparateColumns = true
	t.table.Style().Options.SeparateHeader = true
	cmd.Printf("\n%s\n", t.table.Render())
}
