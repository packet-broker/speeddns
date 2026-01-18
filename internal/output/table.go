package output

import (
	"fmt"
	"io"
	"sort"
	"time"

	"github.com/olekukonko/tablewriter"

	"speeddns/internal/benchmark"
)

// TableFormatter outputs results as ASCII table
type TableFormatter struct {
	writer io.Writer
}

// NewTableFormatter creates a new table formatter
func NewTableFormatter(w io.Writer) *TableFormatter {
	return &TableFormatter{writer: w}
}

// Format outputs results as a formatted table
func (f *TableFormatter) Format(results []benchmark.ResolverResult) error {
	// Filter out results with no successful queries
	validResults := make([]benchmark.ResolverResult, 0, len(results))
	for _, r := range results {
		if r.Successes > 0 {
			validResults = append(validResults, r)
		}
	}

	// Sort by average latency
	sort.Slice(validResults, func(i, j int) bool {
		return validResults[i].Stats.Mean < validResults[j].Stats.Mean
	})

	table := tablewriter.NewWriter(f.writer)
	table.SetHeader([]string{
		"Rank", "Resolver", "IP", "Avg", "Min", "Max",
		"P95", "Success", "Queries",
	})

	// Configure table style
	table.SetBorder(true)
	table.SetRowLine(false)
	table.SetAutoWrapText(false)
	table.SetAlignment(tablewriter.ALIGN_RIGHT)
	table.SetHeaderAlignment(tablewriter.ALIGN_CENTER)
	table.SetColumnAlignment([]int{
		tablewriter.ALIGN_RIGHT,  // Rank
		tablewriter.ALIGN_LEFT,   // Resolver
		tablewriter.ALIGN_LEFT,   // IP
		tablewriter.ALIGN_RIGHT,  // Avg
		tablewriter.ALIGN_RIGHT,  // Min
		tablewriter.ALIGN_RIGHT,  // Max
		tablewriter.ALIGN_RIGHT,  // P95
		tablewriter.ALIGN_RIGHT,  // Success
		tablewriter.ALIGN_RIGHT,  // Queries
	})

	for i, r := range validResults {
		successRate := float64(r.Successes) / float64(r.Queries) * 100

		table.Append([]string{
			fmt.Sprintf("%d", i+1),
			r.Resolver.Name,
			r.Address,
			formatDuration(r.Stats.Mean),
			formatDuration(r.Stats.Min),
			formatDuration(r.Stats.Max),
			formatDuration(r.Stats.P95),
			fmt.Sprintf("%.1f%%", successRate),
			fmt.Sprintf("%d", r.Queries),
		})
	}

	table.Render()

	// Show failed resolvers if any
	failedCount := len(results) - len(validResults)
	if failedCount > 0 {
		fmt.Fprintf(f.writer, "\n%d resolver(s) failed all queries and are not shown.\n", failedCount)
	}

	return nil
}

// formatDuration formats duration for display
func formatDuration(d time.Duration) string {
	if d == 0 {
		return "-"
	}
	if d < time.Millisecond {
		return fmt.Sprintf("%.0fus", float64(d.Microseconds()))
	}
	return fmt.Sprintf("%.2fms", float64(d.Microseconds())/1000)
}
