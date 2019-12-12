Reporter
===
<p>Report is a package that simplifies creating CSV-formatted reports. To create a report, simply create a struct that implements the `Row` interface, populate the report rows, then write the report.<p>

### Features:
* Define custom report structures via the `Row` interface
* Generate report with or without CSV header
* Save report to file, StdOut, output buffer - anything that implements a `io.Writer`

### Examples
#### Basic use
```go
...

// define reporter structure
type ReportRow struct {
	First  string
	Second string
}

// implement interface methods
func (r *ReportRow) Header() []string {
	return []string{"first_thing", "second_thing"}
}

func (r *ReportRow) Slice() []string {
	return []string{r.First, r.Second}
}

func main() {
	reportFile, err := os.Create("my_report.csv")
	if err != nil {
		log.Printf("error creating file - %v", err)
	}
	defer reportFile.Close()

	reportBuilder := reporter.New(reportFile)
	
	// do lots of work to generate and add reporter rows
	reportBuilder.AddRow(&ReportRow{
		First:  "first",
		Second: "second",
	})

	rowsWritten, err := reportBuilder.WriteAll()
	if err != nil {
		log.Fatalf("error writing reporter - %v", err)
	}
	log.Printf("wrote %d lines to %s", rowsWritten, reportFile.Name())

	...
}

...
```

#### Create report without header
```go
...

reportFile, err := os.Create("my_report.csv")
if err != nil {
    log.Printf("error creating file - %v", err)
}
defer reportFile.Close()

reportBuilder := reporter.NewWithoutHeader(reportFile)

...
```

#### Send output to the console
```go
...
reportBuilder := reporter.New(os.Stdout)
...
```

#### Save output in memory
```go
...
buf := bytes.NewBuffer(nil)
reportBuilder := reporter.New(buf)
...
```
