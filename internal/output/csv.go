package output

import (
	"encoding/csv"
	"fmt"
	"io"
	"sort"

	"speeddns/internal/benchmark"
)

// CSVFormatter outputs results as CSV
type CSVFormatter struct {
	writer io.Writer
}

// NewCSVFormatter creates a new CSV formatter
func NewCSVFormatter(w io.Writer) *CSVFormatter {
	return &CSVFormatter{writer: w}
}

// Format outputs results as CSV
func (f *CSVFormatter) Format(results []benchmark.ResolverResult) error {
	// Filter and sort
	validResults := make([]benchmark.ResolverResult, 0, len(results))
	for _, r := range results {
		if r.Successes > 0 {
			validResults = append(validResults, r)
		}
	}

	sort.Slice(validResults, func(i, j int) bool {
		return validResults[i].Stats.Mean < validResults[j].Stats.Mean
	})

	w := csv.NewWriter(f.writer)
	defer w.Flush()

	// Write header
	header := []string{
		"rank", "resolver", "provider", "ip", "avg_ms", "min_ms", "max_ms",
		"median_ms", "p75_ms", "p90_ms", "p95_ms", "p99_ms", "std_dev_ms",
		"success_rate", "queries", "successes", "failures",
	}
	if err := w.Write(header); err != nil {
		return err
	}

	// Write data rows
	for i, r := range validResults {
		successRate := float64(r.Successes) / float64(r.Queries) * 100
		row := []string{
			fmt.Sprintf("%d", i+1),
			r.Resolver.Name,
			r.Resolver.Provider,
			r.Address,
			fmt.Sprintf("%.3f", float64(r.Stats.Mean.Microseconds())/1000),
			fmt.Sprintf("%.3f", float64(r.Stats.Min.Microseconds())/1000),
			fmt.Sprintf("%.3f", float64(r.Stats.Max.Microseconds())/1000),
			fmt.Sprintf("%.3f", float64(r.Stats.Median.Microseconds())/1000),
			fmt.Sprintf("%.3f", float64(r.Stats.P75.Microseconds())/1000),
			fmt.Sprintf("%.3f", float64(r.Stats.P90.Microseconds())/1000),
			fmt.Sprintf("%.3f", float64(r.Stats.P95.Microseconds())/1000),
			fmt.Sprintf("%.3f", float64(r.Stats.P99.Microseconds())/1000),
			fmt.Sprintf("%.3f", float64(r.Stats.StdDev.Microseconds())/1000),
			fmt.Sprintf("%.2f", successRate),
			fmt.Sprintf("%d", r.Queries),
			fmt.Sprintf("%d", r.Successes),
			fmt.Sprintf("%d", r.Failures),
		}
		if err := w.Write(row); err != nil {
			return err
		}
	}

	return nil
}
