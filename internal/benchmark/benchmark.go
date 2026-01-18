package benchmark

import (
	"context"
	"sync"
	"time"

	"speeddns/internal/dns"
	"speeddns/internal/resolver"
	"speeddns/internal/stats"

	mdns "github.com/miekg/dns"
)

// Config holds benchmark configuration
type Config struct {
	Timeout     time.Duration
	Iterations  int
	Concurrency int
	UseTCP      bool
	IncludeIPv6 bool
	Domains     []string
}

// DefaultConfig returns sensible defaults
func DefaultConfig() Config {
	return Config{
		Timeout:     5 * time.Second,
		Iterations:  5,
		Concurrency: 10,
		UseTCP:      false,
		IncludeIPv6: false,
		Domains:     DefaultTestDomains(),
	}
}

// ResolverResult holds aggregated results for a resolver
type ResolverResult struct {
	Resolver  resolver.Resolver `json:"resolver"`
	Address   string            `json:"address"`
	Queries   int               `json:"queries"`
	Successes int               `json:"successes"`
	Failures  int               `json:"failures"`
	RTTs      []time.Duration   `json:"-"`
	Stats     stats.Summary     `json:"stats"`
	Errors    []string          `json:"errors,omitempty"`
}

// Benchmark orchestrates the DNS benchmark tests
type Benchmark struct {
	config    Config
	client    *dns.Client
	resolvers []resolver.Resolver
}

// New creates a new Benchmark instance
func New(config Config, resolvers []resolver.Resolver) *Benchmark {
	return &Benchmark{
		config:    config,
		client:    dns.NewClient(config.Timeout, config.UseTCP),
		resolvers: resolvers,
	}
}

// Progress reports benchmark progress
type Progress struct {
	Resolver string
	Address  string
	Total    int
	Current  int
}

// Run executes the benchmark and returns results
func (b *Benchmark) Run(ctx context.Context, progress chan<- Progress) ([]ResolverResult, error) {
	var wg sync.WaitGroup
	resultsChan := make(chan ResolverResult, len(b.resolvers)*4)

	// Semaphore for concurrency control
	sem := make(chan struct{}, b.config.Concurrency)

	// Count total addresses to test
	total := 0
	for _, res := range b.resolvers {
		total += len(res.AllAddresses(b.config.IncludeIPv6))
	}

	current := 0
	var mu sync.Mutex

	// Test each resolver
	for _, res := range b.resolvers {
		for _, addr := range res.AllAddresses(b.config.IncludeIPv6) {
			wg.Add(1)
			go func(r resolver.Resolver, address string) {
				defer wg.Done()
				sem <- struct{}{}        // acquire
				defer func() { <-sem }() // release

				result := b.testResolver(ctx, r, address)
				resultsChan <- result

				if progress != nil {
					mu.Lock()
					current++
					progress <- Progress{
						Resolver: r.Name,
						Address:  address,
						Total:    total,
						Current:  current,
					}
					mu.Unlock()
				}
			}(res, addr)
		}
	}

	// Close results channel when all done
	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	// Collect results
	var results []ResolverResult
	for result := range resultsChan {
		results = append(results, result)
	}

	return results, nil
}

// testResolver runs all test queries against a single resolver address
func (b *Benchmark) testResolver(ctx context.Context, res resolver.Resolver, addr string) ResolverResult {
	result := ResolverResult{
		Resolver: res,
		Address:  addr,
		RTTs:     make([]time.Duration, 0, b.config.Iterations*len(b.config.Domains)),
	}

	// Early bailout: if first N queries all fail, resolver is likely unreachable
	consecutiveFailures := 0
	const maxConsecutiveFailures = 3

	for i := 0; i < b.config.Iterations; i++ {
		for _, domain := range b.config.Domains {
			select {
			case <-ctx.Done():
				result.Stats = stats.Calculate(result.RTTs)
				return result
			default:
			}

			qr := b.client.Query(ctx, addr, domain, mdns.TypeA)
			result.Queries++

			if qr.Success {
				result.Successes++
				result.RTTs = append(result.RTTs, qr.RTT)
				consecutiveFailures = 0 // reset on success
			} else {
				result.Failures++
				consecutiveFailures++

				// Early bailout: if we've never succeeded and hit max consecutive failures, give up
				if result.Successes == 0 && consecutiveFailures >= maxConsecutiveFailures {
					if len(result.Errors) < 5 {
						result.Errors = append(result.Errors, "early bailout: resolver unreachable")
					}
					result.Stats = stats.Calculate(result.RTTs)
					return result
				}

				if qr.Error != nil {
					// Limit error collection to avoid memory issues
					if len(result.Errors) < 5 {
						result.Errors = append(result.Errors, qr.Error.Error())
					}
				}
			}
		}
	}

	// Calculate statistics
	result.Stats = stats.Calculate(result.RTTs)

	return result
}

// Runner manages benchmark execution with timeout and cancellation
type Runner struct {
	benchmark *Benchmark
	timeout   time.Duration
}

// NewRunner creates a new Runner
func NewRunner(b *Benchmark, timeout time.Duration) *Runner {
	return &Runner{
		benchmark: b,
		timeout:   timeout,
	}
}

// Execute runs the benchmark with overall timeout
func (r *Runner) Execute(progressCallback func(Progress)) ([]ResolverResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	var progress chan Progress
	if progressCallback != nil {
		progress = make(chan Progress, 100)
		go func() {
			for p := range progress {
				progressCallback(p)
			}
		}()
	}

	results, err := r.benchmark.Run(ctx, progress)

	if progress != nil {
		close(progress)
	}

	return results, err
}
