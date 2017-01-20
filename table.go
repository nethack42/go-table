package table

import (
	"bytes"
	"fmt"
	"io"

	"github.com/fatih/color"
)

const (
	Info  = color.FgWhite
	Warn  = color.FgYellow
	Error = color.FgHiRed
)

type Column struct {
	// Value holds the string value of the column
	Value string

	// Attributes can be used instead of Color, Italic, Underline and Bold
	Attributes []color.Attribute

	// Importance holds the values importance.
	Color color.Attribute

	// Italic is set to true if the column should be printed in italic font. This
	// may not work on each and every system
	Italic bool

	// Underline is set to true if the column should be printed with underlines
	Underline bool

	// Bold is set to true if the column should be printed using bold font
	Bold bool

	// RightAlign causes the column to be right aligned
	RightAlign bool
}

type Row []Column

type Table struct {
	// Rows holds a slice of string slices representing rows and columns
	Rows []Row

	Spacing int
}

func (t *Table) Write(w io.Writer) error {
	columnSizes := t.calculateColumnSizes()

	cfg := &RowPrintConfig{
		ColumnSizes: columnSizes,
		Newline:     true,
		Spacing:     t.Spacing,
	}

	for _, row := range t.Rows {
		if err := PrintRow(w, row, cfg); err != nil {
			return err
		}
	}

	return nil
}

func (t *Table) calculateColumnSizes() []int {
	var sizes []int

	cols := t.getMaxColumns()
	if cols == 0 {
		return []int{}
	}

	for i := range t.Rows[0] {
		sizes = append(sizes, t.getColumnSize(i))
	}

	return sizes
}

func (t *Table) getMaxColumns() int {
	var max int
	for _, row := range t.Rows {
		if len(row) > max {
			max = len(row)
		}
	}

	return max
}

func (t *Table) AddRow(r Row) error {
	t.Rows = append(t.Rows, r)

	return nil
}

func (t *Table) getColumnSize(i int) int {
	var max int

	for _, row := range t.Rows {
		if len(row) < i-1 {
			// Empty columns
			continue
		}

		c := row[i]
		if len(c.Value) > max {
			max = len(c.Value)
		}
	}

	return max
}

func getCellPadding(val string, max int) string {
	if len(val) < max {
		padding := ""
		for i := len(val); i < max; i++ {
			padding += " "
		}

		return padding
	}

	return ""
}

type RowPrintConfig struct {
	// Newline specifies if a newline (LF) should be printed
	Newline bool

	ColumnSizes []int

	Spacing int
}

func PrintRow(w io.Writer, row Row, cfg *RowPrintConfig) error {
	var columnSizes []int
	newLine := false
	spacing := ""

	if cfg != nil {
		columnSizes = cfg.ColumnSizes
		newLine = cfg.Newline
		spacing = getCellPadding("", cfg.Spacing)
	}

	// we write into buf rather than into w directly so we can discard it in case
	// of an error
	buf := new(bytes.Buffer)

	// if columnSizes are predefined we can grow our buffer to the estimated
	// amout of memory requrired
	if columnSizes != nil {
		if len(columnSizes) != len(row) {
			return fmt.Errorf("number of column sizes does not match number of columns: sizes=%d cols=%d", len(columnSizes), len(row))
		}
		sum := 0
		for _, size := range columnSizes {
			sum += size
		}

		buf.Grow(sum)
	}

	// print columns
	for i, column := range row {
		padding := ""

		if columnSizes != nil {
			padding = getCellPadding(column.Value, columnSizes[i])
		}

		if err := PrintColumn(buf, column, padding); err != nil {
			return err
		}

		if _, err := buf.Write([]byte(spacing)); err != nil {
			return err
		}
	}

	if newLine {
		if _, err := buf.Write([]byte("\n")); err != nil {
			return err
		}
	}

	// Copy to output provider
	// if this failes, there's nothing we can do ...
	if _, err := w.Write(buf.Bytes()); err != nil {
		return err
	}

	return nil
}

func PrintColumn(w io.Writer, c Column, padding string) error {
	var attr []color.Attribute

	attr = append(attr, c.Color)

	if c.Attributes == nil {
		if c.Bold {
			attr = append(attr, color.Bold)
		}

		if c.Underline {
			attr = append(attr, color.Underline)
		}

		if c.Italic {
			attr = append(attr, color.Italic)
		}

	} else {
		attr = c.Attributes
	}

	printer := color.New(attr...).SprintFunc()

	if c.RightAlign {
		if _, err := w.Write([]byte(padding)); err != nil {
			return err
		}
	}

	if _, err := w.Write([]byte(printer(c.Value))); err != nil {
		return err
	}

	if !c.RightAlign {
		if _, err := w.Write([]byte(padding)); err != nil {
			return err
		}
	}

	return nil
}
