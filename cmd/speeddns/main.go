package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"

	"speeddns/internal/benchmark"
	"speeddns/internal/output"
	"speeddns/internal/resolver"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

// CLI flags
var (
	flagTimeout     time.Duration
	flagIterations  int
	flagConcurrency int
	flagFormat      string
	flagOutput      string
	flagUseTCP      bool
	flagIPv6        bool
	flagQuiet       bool
	flagExtended    bool
	flagResolvers   []string
	flagDomains     []string
	flagListOnly    bool
	flagPrimaryOnly bool
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "speeddns",
		Short: "DNS resolution speed testing tool",
		Long: `speeddns benchmarks DNS resolution speed across multiple public resolvers.

It queries a set of test domains against various DNS servers and reports
timing statistics including min, max, average, and percentiles (p50, p95, p99).

Example usage:
  speeddns                    # Test all built-in resolvers
  speeddns -n 10              # 10 iterations per domain
  speeddns -f json -o out.json # Output JSON to file
  speeddns -r 192.168.1.1     # Add custom resolver
  speeddns --list             # List all built-in resolvers`,
		Version: fmt.Sprintf("%s (commit: %s, built: %s)", version, commit, date),
		RunE:    run,
	}

	// Define flags
	flags := rootCmd.Flags()
	flags.DurationVarP(&flagTimeout, "timeout", "t", 5*time.Second,
		"Timeout for each DNS query")
	flags.IntVarP(&flagIterations, "iterations", "n", 5,
		"Number of query iterations per domain")
	flags.IntVarP(&flagConcurrency, "concurrency", "c", 10,
		"Number of concurrent resolver tests")
	flags.StringVarP(&flagFormat, "format", "f", "table",
		"Output format: table, json, csv")
	flags.StringVarP(&flagOutput, "output", "o", "",
		"Output file (default: stdout)")
	flags.BoolVar(&flagUseTCP, "tcp", false,
		"Use TCP instead of UDP")
	flags.BoolVar(&flagIPv6, "ipv6", false,
		"Include IPv6 resolver addresses")
	flags.BoolVarP(&flagQuiet, "quiet", "q", false,
		"Suppress progress output")
	flags.BoolVar(&flagExtended, "extended", false,
		"Use extended domain list for testing")
	flags.StringSliceVarP(&flagResolvers, "resolver", "r", nil,
		"Additional resolver IPs to test (can be repeated)")
	flags.StringSliceVarP(&flagDomains, "domain", "d", nil,
		"Custom domains to query (can be repeated)")
	flags.BoolVarP(&flagListOnly, "list", "l", false,
		"List built-in resolvers and exit")
	flags.BoolVarP(&flagPrimaryOnly, "primary", "p", false,
		"Only test primary IP of each resolver (faster)")

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func run(cmd *cobra.Command, args []string) error {
	// Handle list-only mode
	if flagListOnly {
		return listResolvers()
	}

	// Handle signals for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		fmt.Fprintln(os.Stderr, "\nInterrupted, stopping...")
		os.Exit(130)
	}()

	// Build resolver list
	resolvers := resolver.BuiltinResolvers()

	// If primary-only mode, reduce to just primary IPs
	if flagPrimaryOnly {
		for i := range resolvers {
			if len(resolvers[i].IPv4) > 1 {
				resolvers[i].IPv4 = resolvers[i].IPv4[:1]
			}
			if len(resolvers[i].IPv6) > 1 {
				resolvers[i].IPv6 = resolvers[i].IPv6[:1]
			}
		}
	}

	// Add custom resolvers if specified
	for _, r := range flagResolvers {
		resolvers = append(resolvers, resolver.Resolver{
			Name:     r,
			Provider: "Custom",
			IPv4:     []string{r},
		})
	}

	// Build configuration
	config := benchmark.Config{
		Timeout:     flagTimeout,
		Iterations:  flagIterations,
		Concurrency: flagConcurrency,
		UseTCP:      flagUseTCP,
		IncludeIPv6: flagIPv6,
	}

	// Set domains
	if len(flagDomains) > 0 {
		config.Domains = flagDomains
	} else if flagExtended {
		config.Domains = benchmark.ExtendedTestDomains()
	} else {
		config.Domains = benchmark.DefaultTestDomains()
	}

	// Print test info
	if !flagQuiet {
		totalAddresses := 0
		for _, r := range resolvers {
			totalAddresses += len(r.AllAddresses(flagIPv6))
		}
		fmt.Fprintf(os.Stderr, "Testing %d resolvers (%d addresses) with %d domains, %d iterations each\n",
			len(resolvers), totalAddresses, len(config.Domains), config.Iterations)
		fmt.Fprintf(os.Stderr, "Total queries per resolver: %d\n\n", len(config.Domains)*config.Iterations)
	}

	// Create and run benchmark
	b := benchmark.New(config, resolvers)
	runner := benchmark.NewRunner(b, time.Hour) // Long timeout for full run

	// Progress callback
	var progressCallback func(benchmark.Progress)
	if !flagQuiet {
		progressCallback = func(p benchmark.Progress) {
			fmt.Fprintf(os.Stderr, "\rTesting resolvers... %d/%d completed", p.Current, p.Total)
		}
	}

	results, err := runner.Execute(progressCallback)
	if err != nil {
		return fmt.Errorf("benchmark failed: %w", err)
	}

	if !flagQuiet {
		fmt.Fprintln(os.Stderr, "\n")
	}

	// Setup output
	var w *os.File = os.Stdout
	if flagOutput != "" {
		f, err := os.Create(flagOutput)
		if err != nil {
			return fmt.Errorf("failed to create output file: %w", err)
		}
		defer f.Close()
		w = f
	}

	// Format and output results
	formatter := output.New(output.Format(flagFormat), w)
	return formatter.Format(results)
}

func listResolvers() error {
	fmt.Println("Built-in DNS Resolvers")
	fmt.Println("======================")
	for _, r := range resolver.BuiltinResolvers() {
		fmt.Printf("\n%s (%s)\n", r.Name, r.Provider)
		fmt.Printf("  Description: %s\n", r.Description)
		fmt.Printf("  IPv4: %v\n", r.IPv4)
		if len(r.IPv6) > 0 {
			fmt.Printf("  IPv6: %v\n", r.IPv6)
		}
		if len(r.Features) > 0 {
			fmt.Printf("  Features: %v\n", r.Features)
		}
	}
	return nil
}
