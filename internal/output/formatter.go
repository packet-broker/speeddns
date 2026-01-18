package output

import (
	"io"

	"speeddns/internal/benchmark"
)

// Format represents output format type
type Format string

const (
	FormatTable Format = "table"
	FormatJSON  Format = "json"
	FormatCSV   Format = "csv"
)

// Formatter defines the output formatting interface
type Formatter interface {
	Format(results []benchmark.ResolverResult) error
}

// New creates a formatter based on format type
func New(format Format, w io.Writer) Formatter {
	switch format {
	case FormatJSON:
		return NewJSONFormatter(w)
	case FormatCSV:
		return NewCSVFormatter(w)
	default:
		return NewTableFormatter(w)
	}
}
