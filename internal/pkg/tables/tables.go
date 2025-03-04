package tables

import (
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/print"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
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

// Sets the title of the table
func (t *Table) SetTitle(title string) {
	t.table.SetTitle(title)

	// prevent title wrapping by setting the width of the first column to the length of the title
	// this is a workaround for a bug in the tables pkg, see https://github.com/jedib0t/go-pretty/issues/135
	t.table.SetColumnConfigs([]table.ColumnConfig{
		{
			Number:   1,
			WidthMin: len(title),
		},
	})
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

// Returns the table rendered
func (t *Table) Render() string {
	t.table.SetStyle(table.StyleLight)

	t.table.Style().Title = table.TitleOptionsBlackOnCyan
	t.table.Style().Title.Align = text.AlignCenter

	t.table.Style().Options.DrawBorder = false
	t.table.Style().Options.SeparateRows = false
	t.table.Style().Options.SeparateColumns = true
	t.table.Style().Options.SeparateHeader = true

	return fmt.Sprintf("\n%s\n\n", t.table.Render())
}

// Displays the table in the command's stdout
func (t *Table) Display(p *print.Printer) error {
	return p.PagerDisplay(t.Render())
}

// Displays multiple tables in the command's stdout
func DisplayTables(p *print.Printer, tables []Table) error {
	renderedTables := ""

	for _, t := range tables {
		renderedTables += t.Render()
	}

	return p.PagerDisplay(renderedTables)
}
