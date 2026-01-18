package benchmark

// DefaultTestDomains returns a balanced set of domains for testing
func DefaultTestDomains() []string {
	return []string{
		// High-traffic global sites (likely cached everywhere)
		"google.com",
		"facebook.com",
		"youtube.com",
		"amazon.com",
		"microsoft.com",

		// Tech/CDN domains
		"cloudflare.com",
		"github.com",
		"netflix.com",
		"apple.com",

		// Regional variety
		"bbc.co.uk",
		"wikipedia.org",

		// Less common (tests resolver's recursive lookup)
		"example.com",
		"iana.org",
	}
}

// ExtendedTestDomains returns a larger set for more thorough testing
func ExtendedTestDomains() []string {
	return []string{
		// Major sites
		"google.com", "facebook.com", "youtube.com", "amazon.com",
		"microsoft.com", "apple.com", "netflix.com", "twitter.com",
		"instagram.com", "linkedin.com", "reddit.com", "github.com",

		// CDNs and infrastructure
		"cloudflare.com", "akamai.com", "fastly.com", "amazonaws.com",

		// News and media
		"bbc.co.uk", "cnn.com", "nytimes.com", "theguardian.com",

		// Tech companies
		"stackoverflow.com", "docker.com", "kubernetes.io", "golang.org",

		// Standards bodies
		"ietf.org", "w3.org", "iana.org", "icann.org",

		// General
		"wikipedia.org", "example.com", "example.org", "example.net",
	}
}
