package reporter

import (
	"bytes"
	"log"
	"os"
	"reflect"
	"testing"
)

func setup(testName string) *os.File {
	f, err := os.Create(testName + ".csv")
	if err != nil {
		log.Fatalf("%T - %v", err, err)
	}
	return f
}

func teardown(f *os.File) {
	if err := f.Close(); err != nil {
		log.Fatalf("%T - %v", err, err)
	}
	if err := os.Remove(f.Name()); err != nil {
		log.Fatalf("%T - %v", err, err)
	}
}

func TestNew(t *testing.T) {
	f := setup(t.Name())
	defer teardown(f)

	r := New(f)
	if !r.includeHeader {
		t.Error("includeHeader false, wanted true")
	}
}

func TestNewWithoutHeader(t *testing.T) {
	f := setup(t.Name())
	defer teardown(f)

	r := NewWithoutHeader(f)
	if r.includeHeader {
		t.Error("includeHeader true, wanted false")
	}
}

func TestReport_Flush(t *testing.T) {
	f := setup(t.Name())
	defer teardown(f)

	r := New(f)
	if err := r.Flush(); err != nil {
		t.Errorf("%T - %v", err, err)
	}
}

func TestReport_AddRow(t *testing.T) {
	r := New(os.Stdout)
	r.AddRow(&testRow{field: "field"})

	if r.Length() != 1 {
		t.Errorf("report.Length()=%d", r.Length())
	}
	if len(r.rows) != r.Length() {
		t.Errorf("len(rows)=%d, report.Length()=%d", len(r.rows), r.Length())
	}
	if row := r.Row(0); row != nil {

	}
}

func TestReport_RemoveRow(t *testing.T) {
	r := New(os.Stdout)
	r.AddRow(&testRow{field: "1"})
	r.AddRow(&testRow{field: "2"})
	if r.Length() != 2 {
		t.Fatalf("report.Length()=%d", r.Length())
	}

	err := r.RemoveRow(0)
	if err != nil {
		t.Fatalf("%T - %v", err, err)
	}
	if r.Length() != 1 {
		t.Errorf("report.Length()=%d", r.Length())
	}
	if row := r.Row(0); row.(*testRow).field != "2" {
		t.Errorf("%#v", row)
	}
}

func TestReport_RemoveRow_outOfBounds(t *testing.T) {
	r := New(os.Stdout)
	r.AddRow(&testRow{field: "1"})
	r.AddRow(&testRow{field: "2"})
	if r.Length() != 2 {
		t.Fatalf("report.Length()=%d", r.Length())
	}

	badIdx := 4
	err := r.RemoveRow(badIdx)
	if err == nil {
		t.Fatalf("expected error - len=%d idx=%d", r.Length(), badIdx)
	}
}

func TestReport_Row(t *testing.T) {
	f := setup(t.Name())
	defer teardown(f)

	r := New(f)
	expectedRow := &testRow{field: "hello"}
	r.AddRow(expectedRow)

	actualRow := r.Row(0)

	if !reflect.DeepEqual(expectedRow, actualRow) {
		t.Errorf("Want %#v, got %#v", expectedRow, actualRow)
	}
}

func TestReport_Row_outOfBounds(t *testing.T) {
	f := setup(t.Name())
	defer teardown(f)

	r := New(f)
	expectedRow := &testRow{field: "hello"}
	r.AddRow(expectedRow)

	badIdx := -1
	actualRow := r.Row(badIdx)
	if actualRow != nil {
		t.Errorf("Want nil, got %#v", actualRow)
	}
}

func TestReport_WriteAll(t *testing.T) {
	buf := bytes.NewBuffer([]byte{})
	r := New(buf)

	r.AddRow(&testRow{field: "1"})
	r.AddRow(&testRow{field: "2"})
	if r.Length() != 2 {
		t.Fatalf("report.Length()=%d", r.Length())
	}

	rowsWritten, err := r.WriteAll()
	if err != nil {
		t.Fatalf("%T - %v", err, err)
	}
	if rowsWritten != r.Length() {
		t.Errorf("rowsWritten=%d, report.Length()=%d", rowsWritten, r.Length())
	}
	expected := "field\n1\n2\n"
	if buf.String() != expected {
		t.Errorf("Want %q, got %q", expected, buf.String())
	}
}

func TestReport_WriteAll_noHeader(t *testing.T) {
	buf := bytes.NewBuffer([]byte{})
	r := NewWithoutHeader(buf)

	r.AddRow(&testRow{field: "1"})
	r.AddRow(&testRow{field: "2"})
	if r.Length() != 2 {
		t.Fatalf("report.Length()=%d", r.Length())
	}

	rowsWritten, err := r.WriteAll()
	if err != nil {
		t.Fatalf("%T - %v", err, err)
	}
	if rowsWritten != r.Length() {
		t.Errorf("rowsWritten=%d, report.Length()=%d", rowsWritten, r.Length())
	}
	expected := "1\n2\n"
	if buf.String() != expected {
		t.Errorf("Want %q, got %q", expected, buf.String())
	}
}

func TestReport_WriteAll_fileClosed(t *testing.T) {
	f := setup(t.Name())

	r := New(f)
	r.AddRow(&testRow{field: "1"})
	r.AddRow(&testRow{field: "2"})
	if r.Length() != 2 {
		t.Fatalf("report.Length()=%d", r.Length())
	}

	teardown(f)

	rowsWritten, err := r.WriteAll()
	if err == nil || rowsWritten > 0 {
		t.Errorf("err=%v, rowsWritten=%d", err, rowsWritten)
	}
}

type testRow struct {
	field string
}

func (r *testRow) Header() []string {
	return []string{"field"}
}

func (r *testRow) Marshal() []string {
	return []string{r.field}
}
