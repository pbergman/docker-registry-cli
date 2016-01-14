package helpers

import (
	"fmt"
	"strconv"
	"strings"
	"io"
	"os"
)

type Table struct {
	Header []string
	Rows   [][]string
}

func NewTable(header ...string) *Table {
	return &Table{
		Header: header,
		Rows:	make([][]string, 0),
	}
}

func (t *Table) GetColumnLen() int {
	return len(t.Header)
}

func (t *Table) AddRow(data ...interface{}) {

	row := make([]string, len(data))

	for i, d := range data {
		switch d.(type) {
		case int:
			row[i] = strconv.Itoa(d.(int))
		case string:
			row[i] = d.(string)
		default:
			row[i] = fmt.Sprintf("%#v", d)
		}
	}

	t.Rows = append(t.Rows, row)
}

func (t *Table) pad(string, padding string, width int) string {
	if left := width - len(string); left > 0 {
		return string + strings.Repeat(padding, left)
	}
	return string
}
// getSizes returns a slice with int`s representing
// the max size of every column:
// [3]int{5,8,7} so first column is 5 characters
// second 8 and third 7 characters width
func (t *Table) getSizes() []int {
	sizes := make([]int, t.GetColumnLen())
	for i, name := range t.Header {
		sizes[i] = len(name)
	}
	for _, row := range t.Rows {
		for i, column := range row {
			if sizes[i] < len(column) {
				sizes[i] = len(column)
			}
		}
	}
	return sizes
}

func (t *Table) Print() {
	sizes := t.getSizes()
	border := make([]string, 0)
	for i := 0; i < len(sizes); i++ {
		border = append(border, "=")
	}
	writeLine := func(padding string, data []string) {
		for i, str := range data {
			io.WriteString(os.Stdout, t.pad(str, padding, sizes[i]))
			io.WriteString(os.Stdout, " ")
		}
		io.WriteString(os.Stdout, "\n")
	}
	writeLine("=", border)
	writeLine(" ", t.Header)
	writeLine("=", border)
	for _, row := range t.Rows {
		writeLine(" ", row)
	}
	writeLine("=", border)
}