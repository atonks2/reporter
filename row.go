package reporter

import (
	"reflect"
)

// Row provides the methods needed to convert a struct to a CSV row
type Row interface {
	// Header should return a list of the column names to be included in the CSV report
	Header() []string
	// Slice should return a list of the data points to be written to the CSV report
	Slice() []string
}

// TimeFormatString defaults to "2006-01-02T15:04:05Z07:00", but any valid time.Time format string can be used
var timeFormatString = "2006-01-02T15:04:05Z"

// SetDateTimeFormat updates the default format string to the supplied string
// Valid formats: https://golang.org/src/time/format.go
func SetDateTimeFormat(format string) {
	timeFormatString = format
}

// CreateHeader uses the "csv" struct tags to build the CSV header
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
