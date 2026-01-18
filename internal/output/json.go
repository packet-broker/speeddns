package output

import (
	"encoding/json"
	"io"
	"sort"

	"speeddns/internal/benchmark"
)

// JSONFormatter outputs results as JSON
type JSONFormatter struct {
	writer io.Writer
}

// NewJSONFormatter creates a new JSON formatter
func NewJSONFormatter(w io.Writer) *JSONFormatter {
	return &JSONFormatter{writer: w}
}

// JSONResult is a JSON-friendly result structure
type JSONResult struct {
	Rank        int     `json:"rank"`
	Name        string  `json:"name"`
	Provider    string  `json:"provider"`
	Address     string  `json:"address"`
	AvgMs       float64 `json:"avg_ms"`
	MinMs       float64 `json:"min_ms"`
	MaxMs       float64 `json:"max_ms"`
	MedianMs    float64 `json:"median_ms"`
	P75Ms       float64 `json:"p75_ms"`
	P90Ms       float64 `json:"p90_ms"`
	P95Ms       float64 `json:"p95_ms"`
	P99Ms       float64 `json:"p99_ms"`
	StdDevMs    float64 `json:"std_dev_ms"`
	SuccessRate float64 `json:"success_rate"`
	Queries     int     `json:"queries"`
	Successes   int     `json:"successes"`
	Failures    int     `json:"failures"`
}

// JSONOutput wraps the results with metadata
type JSONOutput struct {
	Results []JSONResult `json:"results"`
	Summary struct {
		TotalResolvers int `json:"total_resolvers"`
		SuccessfulOnly int `json:"successful_only"`
	} `json:"summary"`
}

// Format outputs results as JSON
func (f *JSONFormatter) Format(results []benchmark.ResolverResult) error {
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

	output := JSONOutput{}
	output.Summary.TotalResolvers = len(results)
	output.Summary.SuccessfulOnly = len(validResults)

	for i, r := range validResults {
		successRate := float64(r.Successes) / float64(r.Queries) * 100

		output.Results = append(output.Results, JSONResult{
			Rank:        i + 1,
			Name:        r.Resolver.Name,
			Provider:    r.Resolver.Provider,
			Address:     r.Address,
			AvgMs:       float64(r.Stats.Mean.Microseconds()) / 1000,
			MinMs:       float64(r.Stats.Min.Microseconds()) / 1000,
			MaxMs:       float64(r.Stats.Max.Microseconds()) / 1000,
			MedianMs:    float64(r.Stats.Median.Microseconds()) / 1000,
			P75Ms:       float64(r.Stats.P75.Microseconds()) / 1000,
			P90Ms:       float64(r.Stats.P90.Microseconds()) / 1000,
			P95Ms:       float64(r.Stats.P95.Microseconds()) / 1000,
			P99Ms:       float64(r.Stats.P99.Microseconds()) / 1000,
			StdDevMs:    float64(r.Stats.StdDev.Microseconds()) / 1000,
			SuccessRate: successRate,
			Queries:     r.Queries,
			Successes:   r.Successes,
			Failures:    r.Failures,
		})
	}

	encoder := json.NewEncoder(f.writer)
	encoder.SetIndent("", "  ")
	return encoder.Encode(output)
}
