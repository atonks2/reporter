package reporter

import (
	"encoding/csv"
	"fmt"
	"io"
)

// Reporter provides functionality for generating CSV-formatted reports
// using any struct that implements the Row interface.
type Reporter struct {
	out           *io.Writer
	w             *csv.Writer
	rows          []*Row
	includeHeader bool
	headerWritten bool
}

// New returns a pointer to a Reporter.
func New(out io.Writer) *Reporter {
	return &Reporter{
		out:           &out,
		w:             csv.NewWriter(out),
		rows:          make([]*Row, 0),
		includeHeader: true,
		headerWritten: false,
	}
}

// NewWithoutHeader returns a pointer to a Reporter that will not include headers in the output
func NewWithoutHeader(out io.Writer) *Reporter {
	r := New(out)
	r.includeHeader = false
	return r
}

// AddRow appends the provided row to Reporter's []*Rows slice
func (r *Reporter) AddRow(row Row) {
	r.rows = append(r.rows, &row)
}

// RemoveRow removes the row at the specified index
func (r *Reporter) RemoveRow(idx int) error {
	if idx < 0 || idx > len(r.rows) {
		return fmt.Errorf("index (%d) out of bounds", idx)
	}
	r.rows[idx] = r.rows[len(r.rows)-1]
	r.rows = r.rows[:len(r.rows)-1]
	return nil
}

// Row retrieves the row at the specified idx, or nil if the index is out of bounds
func (r *Reporter) Row(idx int) Row {
	if idx < 0 || idx > len(r.rows) {
		return nil
	}
	return *r.rows[idx]
}

// Length returns the number of rows currently saved in the Reporter
func (r *Reporter) Length() int {
	return len(r.rows)
}

// WriteAll writes all of the rows saved in the Reporter along with the header, if applicable
func (r *Reporter) WriteAll() (int, error) {
	if r.includeHeader && !r.headerWritten {
		if err := r.writeHeader(); err != nil {
			return 0, err
		}
	}

	numRows := 0
	for i, rowPtr := range r.rows {
		row := *rowPtr
		if err := r.w.Write(row.Slice()); err != nil {
			return 0, fmt.Errorf("failed to write row %d - %v", i, err)
		}
		if err := r.Flush(); err != nil {
			return 0, err
		}
		numRows++
	}
	return numRows, nil
}

// Flush calls the same method in Reporter's underlying csv.Writer, then checks for and returns any errors
func (r *Reporter) Flush() error {
	r.w.Flush()
	return r.w.Error()
}

func (r *Reporter) writeHeader() error {
	if len(r.rows) == 0 {
		return nil
	}
	if !r.headerWritten {
		row := *r.rows[0]
		h := row.Header()
		if err := r.w.Write(h); err != nil {
			return err
		}
		if err := r.Flush(); err != nil {
			return err
		}
		r.headerWritten = true
	}
	return nil
}
