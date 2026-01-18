package dns

import (
	"context"
	"time"

	"github.com/miekg/dns"
)

// QueryResult holds the result of a single DNS query
type QueryResult struct {
	Resolver     string
	Domain       string
	QueryType    uint16
	RTT          time.Duration
	Success      bool
	Error        error
	ResponseCode int
	AnswerCount  int
}

// Client wraps the miekg/dns client with our configuration
type Client struct {
	client  *dns.Client
	timeout time.Duration
}

// NewClient creates a new DNS client with specified timeout
func NewClient(timeout time.Duration, useTCP bool) *Client {
	c := &dns.Client{
		Timeout: timeout,
	}
	if useTCP {
		c.Net = "tcp"
	}
	return &Client{
		client:  c,
		timeout: timeout,
	}
}

// Query performs a DNS query and returns timing information
func (c *Client) Query(ctx context.Context, server, domain string, qtype uint16) QueryResult {
	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn(domain), qtype)
	m.RecursionDesired = true

	result := QueryResult{
		Resolver:  server,
		Domain:    domain,
		QueryType: qtype,
	}

	r, rtt, err := c.client.ExchangeContext(ctx, m, server+":53")
	if err != nil {
		result.Error = err
		result.Success = false
		return result
	}

	result.RTT = rtt
	result.Success = r.Rcode == dns.RcodeSuccess
	result.ResponseCode = r.Rcode
	result.AnswerCount = len(r.Answer)

	return result
}
