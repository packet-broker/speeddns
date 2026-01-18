# SpeedDNS

A fast DNS resolver benchmarking tool for Linux. Test and compare DNS resolution speeds across 20+ public resolvers.

## Quick Start

Download the binary and run - no installation required:

```bash
# Download (choose your architecture)
wget https://github.com/brownbananaz/speeddns/releases/latest/download/speeddns-linux-amd64
chmod +x speeddns-linux-amd64

# Run it
./speeddns-linux-amd64
```

Or move to your PATH for easier access:

```bash
sudo mv speeddns-linux-amd64 /usr/local/bin/speeddns
speeddns
```

## Features

- Tests 20+ public DNS resolvers (Cloudflare, Google, Quad9, OpenDNS, AdGuard, and more)
- Measures latency statistics: min, max, average, and percentiles (P50, P95, P99)
- Multiple output formats: table, JSON, CSV
- Parallel testing for fast results
- Add custom resolvers
- IPv4 and IPv6 support

## Usage

```bash
# Basic test - all resolvers
speeddns

# Faster test - primary IPs only
speeddns -p

# More iterations for accuracy
speeddns -n 10

# Output to JSON
speeddns -f json -o results.json

# Output to CSV
speeddns -f csv > results.csv

# Add a custom resolver
speeddns -r 192.168.1.1

# Test specific domains
speeddns -d example.com -d mysite.org

# Include IPv6 addresses
speeddns --ipv6

# List all built-in resolvers
speeddns --list
```

## Options

| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--timeout` | `-t` | Timeout per query | 5s |
| `--iterations` | `-n` | Queries per domain | 5 |
| `--concurrency` | `-c` | Parallel tests | 10 |
| `--format` | `-f` | Output: table/json/csv | table |
| `--output` | `-o` | Output file | stdout |
| `--primary` | `-p` | Primary IP only (faster) | false |
| `--tcp` | | Use TCP instead of UDP | false |
| `--ipv6` | | Include IPv6 addresses | false |
| `--quiet` | `-q` | Suppress progress | false |
| `--resolver` | `-r` | Add custom resolver | - |
| `--domain` | `-d` | Custom test domain | - |
| `--list` | `-l` | List resolvers | - |
| `--extended` | | Extended domain list | false |

## Sample Output

```
+------+---------------------+----------------+----------+----------+----------+---------+---------+
| RANK |      RESOLVER       |       IP       |   AVG    |   MIN    |   MAX    | SUCCESS | QUERIES |
+------+---------------------+----------------+----------+----------+----------+---------+---------+
|    1 | Cloudflare          | 1.1.1.1        | 12.45ms  |  8.23ms  | 45.12ms  |  100.0% |      65 |
|    2 | Google              | 8.8.8.8        | 15.67ms  |  9.45ms  | 52.34ms  |  100.0% |      65 |
|    3 | Quad9               | 9.9.9.9        | 18.23ms  | 10.12ms  | 78.45ms  |   98.5% |      65 |
|    4 | OpenDNS             | 208.67.222.222 | 22.34ms  | 12.45ms  | 89.12ms  |  100.0% |      65 |
+------+---------------------+----------------+----------+----------+----------+---------+---------+
```

## Built-in Resolvers

| Resolver | Primary IP | Description |
|----------|-----------|-------------|
| Cloudflare | 1.1.1.1 | Privacy-focused, fastest public DNS |
| Google | 8.8.8.8 | High availability, global network |
| Quad9 | 9.9.9.9 | Security-focused, malware blocking |
| OpenDNS | 208.67.222.222 | Phishing protection |
| AdGuard | 94.140.14.14 | Ad-blocking DNS |
| Cloudflare-Malware | 1.1.1.2 | Malware blocking |
| Cloudflare-Family | 1.1.1.3 | Family-safe filtering |
| CleanBrowsing | 185.228.168.9 | Security filter |
| NextDNS | 45.90.28.0 | Customizable DNS |
| Control-D | 76.76.2.0 | Privacy DNS |
| Mullvad | 194.242.2.2 | VPN provider DNS |
| dns0.eu | 193.110.81.0 | European privacy DNS |

Run `speeddns --list` to see all resolvers with their full details.

## Building from Source

Requires Go 1.21+:

```bash
git clone https://github.com/brownbananaz/speeddns.git
cd speeddns
make build
./speeddns
```

Cross-compile for distribution:

```bash
make release          # Build for Linux amd64 and arm64
make dist             # Create .tar.gz archives
```

## License

MIT
