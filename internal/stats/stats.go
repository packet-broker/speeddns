package stats

import (
	"math"
	"sort"
	"time"
)

// Summary holds calculated statistics
type Summary struct {
	Count  int           `json:"count"`
	Min    time.Duration `json:"min"`
	Max    time.Duration `json:"max"`
	Mean   time.Duration `json:"mean"`
	Median time.Duration `json:"median"`
	StdDev time.Duration `json:"std_dev"`
	P50    time.Duration `json:"p50"`
	P75    time.Duration `json:"p75"`
	P90    time.Duration `json:"p90"`
	P95    time.Duration `json:"p95"`
	P99    time.Duration `json:"p99"`
}

// Calculate computes all statistics from RTT samples
func Calculate(rtts []time.Duration) Summary {
	if len(rtts) == 0 {
		return Summary{}
	}

	// Sort for percentile calculations
	sorted := make([]time.Duration, len(rtts))
	copy(sorted, rtts)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i] < sorted[j]
	})

	sum := time.Duration(0)
	for _, rtt := range sorted {
		sum += rtt
	}

	mean := sum / time.Duration(len(sorted))

	// Calculate standard deviation
	var sumSquares float64
	for _, rtt := range sorted {
		diff := float64(rtt - mean)
		sumSquares += diff * diff
	}
	variance := sumSquares / float64(len(sorted))
	stdDev := time.Duration(math.Sqrt(variance))

	return Summary{
		Count:  len(sorted),
		Min:    sorted[0],
		Max:    sorted[len(sorted)-1],
		Mean:   mean,
		Median: percentile(sorted, 50),
		StdDev: stdDev,
		P50:    percentile(sorted, 50),
		P75:    percentile(sorted, 75),
		P90:    percentile(sorted, 90),
		P95:    percentile(sorted, 95),
		P99:    percentile(sorted, 99),
	}
}

// percentile calculates the p-th percentile of a sorted slice
func percentile(sorted []time.Duration, p float64) time.Duration {
	if len(sorted) == 0 {
		return 0
	}
	if len(sorted) == 1 {
		return sorted[0]
	}

	// Linear interpolation method
	rank := (p / 100.0) * float64(len(sorted)-1)
	lower := int(rank)
	upper := lower + 1
	if upper >= len(sorted) {
		upper = len(sorted) - 1
	}

	fraction := rank - float64(lower)
	return sorted[lower] + time.Duration(fraction*float64(sorted[upper]-sorted[lower]))
}
