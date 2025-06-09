package githubip

import (
	"net/netip"
)

//go:generate go run update_ips.go

// IsGithubIP returns true if the ip address falls within one of the known
// GitHub ip ranges.
func IsGithubIP(ip netip.Addr) bool {
	r := Range(ip)
	return r != nil
}

// Range returns the ip range and metadata an address falls within.
// If the IP is not a GitHub IP address it returns nil
func Range(ip netip.Addr) *IPRange {
	_, r, ok := cidrTbl.Lookup(ip)
	if ok {
		return &r
	}
	return nil
}

type IPRange struct {
	Prefix  netip.Prefix
	Service string
}
