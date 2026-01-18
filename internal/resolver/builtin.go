package resolver

// Resolver represents a DNS resolver with metadata
type Resolver struct {
	Name        string   `json:"name"`
	Provider    string   `json:"provider"`
	IPv4        []string `json:"ipv4"`
	IPv6        []string `json:"ipv6,omitempty"`
	Description string   `json:"description"`
	Features    []string `json:"features,omitempty"`
}

// PrimaryAddress returns the primary IPv4 address
func (r Resolver) PrimaryAddress() string {
	if len(r.IPv4) > 0 {
		return r.IPv4[0]
	}
	if len(r.IPv6) > 0 {
		return r.IPv6[0]
	}
	return ""
}

// AllAddresses returns all IP addresses (IPv4 and optionally IPv6)
func (r Resolver) AllAddresses(includeIPv6 bool) []string {
	addrs := make([]string, 0, len(r.IPv4)+len(r.IPv6))
	addrs = append(addrs, r.IPv4...)
	if includeIPv6 {
		addrs = append(addrs, r.IPv6...)
	}
	return addrs
}

// BuiltinResolvers returns the default list of public DNS resolvers
func BuiltinResolvers() []Resolver {
	return []Resolver{
		{
			Name:        "Cloudflare",
			Provider:    "Cloudflare, Inc.",
			IPv4:        []string{"1.1.1.1", "1.0.0.1"},
			IPv6:        []string{"2606:4700:4700::1111", "2606:4700:4700::1001"},
			Description: "Privacy-focused, fastest public DNS",
			Features:    []string{"DoH", "DoT", "DNSSEC"},
		},
		{
			Name:        "Google",
			Provider:    "Google LLC",
			IPv4:        []string{"8.8.8.8", "8.8.4.4"},
			IPv6:        []string{"2001:4860:4860::8888", "2001:4860:4860::8844"},
			Description: "Google Public DNS, high availability",
			Features:    []string{"DoH", "DoT", "DNSSEC", "DNS64"},
		},
		{
			Name:        "Quad9",
			Provider:    "Quad9 Foundation",
			IPv4:        []string{"9.9.9.9", "149.112.112.112"},
			IPv6:        []string{"2620:fe::fe", "2620:fe::9"},
			Description: "Security-focused with malware blocking",
			Features:    []string{"DoH", "DoT", "DNSSEC", "Threat-blocking"},
		},
		{
			Name:        "Quad9-Unfiltered",
			Provider:    "Quad9 Foundation",
			IPv4:        []string{"9.9.9.10", "149.112.112.10"},
			IPv6:        []string{"2620:fe::10", "2620:fe::fe:10"},
			Description: "Quad9 without security filtering",
			Features:    []string{"DoH", "DoT", "DNSSEC"},
		},
		{
			Name:        "OpenDNS",
			Provider:    "Cisco Systems",
			IPv4:        []string{"208.67.222.222", "208.67.220.220"},
			IPv6:        []string{"2620:119:35::35", "2620:119:53::53"},
			Description: "OpenDNS with security features",
			Features:    []string{"DoH", "DNSSEC", "Phishing-protection"},
		},
		{
			Name:        "OpenDNS-FamilyShield",
			Provider:    "Cisco Systems",
			IPv4:        []string{"208.67.222.123", "208.67.220.123"},
			Description: "OpenDNS with adult content filtering",
			Features:    []string{"Content-filtering"},
		},
		{
			Name:        "Cloudflare-Malware",
			Provider:    "Cloudflare, Inc.",
			IPv4:        []string{"1.1.1.2", "1.0.0.2"},
			IPv6:        []string{"2606:4700:4700::1112", "2606:4700:4700::1002"},
			Description: "Cloudflare with malware blocking",
			Features:    []string{"DoH", "DoT", "Malware-blocking"},
		},
		{
			Name:        "Cloudflare-Family",
			Provider:    "Cloudflare, Inc.",
			IPv4:        []string{"1.1.1.3", "1.0.0.3"},
			IPv6:        []string{"2606:4700:4700::1113", "2606:4700:4700::1003"},
			Description: "Cloudflare with malware and adult content blocking",
			Features:    []string{"DoH", "DoT", "Content-filtering"},
		},
		{
			Name:        "AdGuard",
			Provider:    "AdGuard Software Ltd.",
			IPv4:        []string{"94.140.14.14", "94.140.15.15"},
			IPv6:        []string{"2a10:50c0::ad1:ff", "2a10:50c0::ad2:ff"},
			Description: "Ad-blocking DNS service",
			Features:    []string{"DoH", "DoT", "Ad-blocking"},
		},
		{
			Name:        "AdGuard-Family",
			Provider:    "AdGuard Software Ltd.",
			IPv4:        []string{"94.140.14.15", "94.140.15.16"},
			IPv6:        []string{"2a10:50c0::bad1:ff", "2a10:50c0::bad2:ff"},
			Description: "AdGuard with family protection",
			Features:    []string{"DoH", "DoT", "Ad-blocking", "Content-filtering"},
		},
		{
			Name:        "CleanBrowsing-Security",
			Provider:    "CleanBrowsing",
			IPv4:        []string{"185.228.168.9", "185.228.169.9"},
			IPv6:        []string{"2a0d:2a00:1::2", "2a0d:2a00:2::2"},
			Description: "Security filter, blocks malware",
			Features:    []string{"DoH", "DoT", "Security"},
		},
		{
			Name:        "CleanBrowsing-Family",
			Provider:    "CleanBrowsing",
			IPv4:        []string{"185.228.168.168", "185.228.169.168"},
			IPv6:        []string{"2a0d:2a00:1::", "2a0d:2a00:2::"},
			Description: "Family filter with content blocking",
			Features:    []string{"DoH", "DoT", "Content-filtering"},
		},
		{
			Name:        "Comodo-Secure",
			Provider:    "Comodo Group",
			IPv4:        []string{"8.26.56.26", "8.20.247.20"},
			Description: "Comodo Secure DNS",
			Features:    []string{"Malware-blocking"},
		},
		{
			Name:        "Neustar-Reliability",
			Provider:    "Neustar Inc.",
			IPv4:        []string{"64.6.64.6", "64.6.65.6"},
			Description: "Neustar UltraDNS Public",
			Features:    []string{"DNSSEC"},
		},
		{
			Name:        "Neustar-Security",
			Provider:    "Neustar Inc.",
			IPv4:        []string{"156.154.70.2", "156.154.71.2"},
			Description: "Neustar with threat protection",
			Features:    []string{"DNSSEC", "Threat-blocking"},
		},
		{
			Name:        "Level3",
			Provider:    "Level 3 Communications",
			IPv4:        []string{"4.2.2.1", "4.2.2.2"},
			Description: "Level3 Public DNS",
		},
		{
			Name:        "dns0.eu",
			Provider:    "dns0.eu",
			IPv4:        []string{"193.110.81.0", "185.253.5.0"},
			IPv6:        []string{"2a0f:fc80::", "2a0f:fc81::"},
			Description: "European privacy-focused DNS",
			Features:    []string{"DoH", "DoT", "DNSSEC"},
		},
		{
			Name:        "Mullvad",
			Provider:    "Mullvad VPN AB",
			IPv4:        []string{"194.242.2.2"},
			IPv6:        []string{"2a07:e340::2"},
			Description: "Mullvad public DNS",
			Features:    []string{"DoH", "DoT", "Ad-blocking"},
		},
		{
			Name:        "Control-D",
			Provider:    "Control D",
			IPv4:        []string{"76.76.2.0", "76.76.10.0"},
			IPv6:        []string{"2606:1a40::", "2606:1a40:1::"},
			Description: "Control D public DNS",
			Features:    []string{"DoH", "DoT"},
		},
		{
			Name:        "NextDNS",
			Provider:    "NextDNS Inc.",
			IPv4:        []string{"45.90.28.0", "45.90.30.0"},
			IPv6:        []string{"2a07:a8c0::", "2a07:a8c1::"},
			Description: "NextDNS public resolver",
			Features:    []string{"DoH", "DoT", "Customizable"},
		},
		{
			Name:        "Verisign",
			Provider:    "Verisign Inc.",
			IPv4:        []string{"64.6.64.6", "64.6.65.6"},
			IPv6:        []string{"2620:74:1b::1:1", "2620:74:1c::2:2"},
			Description: "Verisign Public DNS",
			Features:    []string{"DNSSEC"},
		},
	}
}
