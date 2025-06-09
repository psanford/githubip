package githubip

import (
	"fmt"
	"net/netip"
	"testing"
)

func TestIsGithubIP(t *testing.T) {
	// These tests will work once we generate ips.gen.go
	// For now, we'll use placeholder IPs that should be GitHub IPs
	githubIPs := []netip.Addr{
		netip.MustParseAddr("192.30.252.1"),
		netip.MustParseAddr("185.199.108.1"),
		netip.MustParseAddr("140.82.112.1"),
		netip.MustParseAddr("143.55.64.1"),
		netip.MustParseAddr("2a0a:a440::1"),
		netip.MustParseAddr("2606:50c0::1"),
	}

	for _, addr := range githubIPs {
		if !IsGithubIP(addr) {
			t.Errorf("Expected %s to match GitHub ip but did not", addr)
		}
	}

	nonGithubIPs := []netip.Addr{
		netip.MustParseAddr("127.0.0.1"),
		netip.MustParseAddr("10.0.0.1"),
		netip.MustParseAddr("8.8.8.8"),
		netip.MustParseAddr("1.1.1.1"),
		netip.MustParseAddr("2606:4700:4700::1111"),
	}
	for _, addr := range nonGithubIPs {
		if IsGithubIP(addr) {
			t.Errorf("%s is not a GitHub ip address, but it matched", addr)
		}
	}
}

func BenchmarkLookup(b *testing.B) {
	tests := []struct {
		ip       string
		isGithub bool
	}{
		{"192.30.252.1", true},
		{"185.199.108.1", true},
		{"140.82.112.1", true},
		{"143.55.64.1", true},
		{"2a0a:a440::1", true},
		{"2606:50c0::1", true},
		{"127.0.0.1", false},
		{"10.0.0.1", false},
		{"8.8.8.8", false},
		{"1.1.1.1", false},
		{"2606:4700:4700::1111", false},
	}

	for _, tc := range tests {
		ip := netip.MustParseAddr(tc.ip)
		b.Run(tc.ip, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				isGithub := IsGithubIP(ip)

				if isGithub != tc.isGithub {
					b.Fatalf("%s got isGithub=%t expected=%t", tc.ip, isGithub, tc.isGithub)
				}
			}
		})
	}
}

func ExampleRange() {
	ip := netip.MustParseAddr("192.30.252.1")
	r := Range(ip)
	fmt.Println(r.Prefix)
	fmt.Println(r.Service)
	// Output:
	// 192.30.252.0/22
	// api
}

func ExampleIsGithubIP() {
	ips := []netip.Addr{
		netip.MustParseAddr("192.30.252.1"),
		netip.MustParseAddr("127.0.0.1"),
	}
	for _, ip := range ips {
		if IsGithubIP(ip) {
			fmt.Printf("%s is GitHub\n", ip)
		} else {
			fmt.Printf("%s is NOT GitHub\n", ip)
		}
	}
	// Output:
	// 192.30.252.1 is GitHub
	// 127.0.0.1 is NOT GitHub
}
