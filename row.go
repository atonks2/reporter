package reporter

import (
	"reflect"
)

// Row provides the methods needed to convert a struct to a CSV row
//
// Example (using experimental reporter.CreateHeader() and reporter.MarshalCSV()):
//
//	type MyData struct {
//		Field1 string    `csv:"field_1"`
//		Field2 time.Time `csv:"field_2"`
//	}
//
//	func (d *MyData) Header() []string {
//		return reporter.CreateHeader(d)
//	}
//
//	func (d *MyData) Marshal() []string {
//		return reporter.MarshalCSV(d)
//	}
type Row interface {
	// Header should return a list of the column names to be included in the CSV report
	Header() []string
	// Marshal should return a list of the data points to be written to the CSV report
	Marshal() []string
}

// TimeFormatString defaults to "2006-01-02T15:04:05Z", but any valid time.Time format string can be used
var timeFormatString = "2006-01-02T15:04:05Z"

// SetDateTimeFormat updates the default format string to the supplied string.
// Valid formats: https://golang.org/src/time/format.go
// This is only used by MarshalCSV.
func SetDateTimeFormat(format string) {
	timeFormatString = format
}

// CreateHeader uses the "csv" struct tags to build the CSV header
// v should be a struct with "csv" tags.
// Example:
//	type MyData struct {
//		Field1   string    `csv:"field_1"`
//		IgnoreMe chan int  `csv:"-"`
//		Field2   time.Time `csv:"field_2"`
//	}
//
// CreateHeader(MyData{}) will output []string{"field_1", "field_2"}, ignoring the "-" tag.
//
// Experimental
func CreateHeader(v interface{}) []string {
	ref := createReflection(v)
	var header []string
	for i := 0; i < ref.Value.NumField(); i++ {
		headerVal := ref.Type.Field(i).Tag.Get("csv")
		if headerVal != "-" {
			header = append(header, headerVal)
		}
	}
	return header
}

// MarshalCSV returns the struct values as a slice of strings.
// Fields with the tag `csv:"-"` will be ignored.
//
// Experimental
func MarshalCSV(v interface{}) []string {
	ref := createReflection(v)
	if ref.Kind != reflect.Struct {
		return []string{ref.getUnderlyingData()}
	}
	var csvRow []string
	for i := 0; i < ref.Value.NumField(); i++ { // loop through struct fields
		if ref.Value.Type().Field(i).Tag.Get("csv") == "-" {
			continue
		}
		fieldRef := ref.getStructFieldAtIdx(i)
		str := fieldRef.tryCommonTypes()
		if str != "" {
			csvRow = append(csvRow, str)
			continue
		}
		csvRow = append(csvRow, fieldRef.getUnderlyingData())
	}
	return csvRow
}
