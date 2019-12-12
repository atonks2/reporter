package reporter

// Row provides the methods needed to convert a struct to a CSV row
type Row interface {
	// Header should return a list of the column names to be included in the CSV report
	Header() []string
	// Slice should return a list of the data points to be written to the CSV report
	Slice() []string
}
